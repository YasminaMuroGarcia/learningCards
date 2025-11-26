package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"learning-cards/internal/handlers"
	"learning-cards/internal/models"
	"learning-cards/internal/repository"
	"learning-cards/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTest creates an in-memory sqlite DB, migrates models, seeds some data and
// returns a handler wired with real repository/service and the DB for assertions.
func setupTest(t *testing.T) (*handlers.UserWordHandler, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite DB: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Word{}, &models.UserWord{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	// Create repository/service/handler
	repo := repository.NewUserWordRepository(db)
	svc := services.NewUserWordService(repo)
	handler := handlers.NewUserWordHandler(svc)

	return handler, db
}

// seedData inserts a few words and one user_word (for the first word) into the DB.
// Returns the created words slice.
func seedData(t *testing.T, db *gorm.DB) []models.Word {
	t.Helper()

	words := []models.Word{
		{Word: "cat", Translation: "gato", Category: "animals", CreatedAt: time.Now()},
		{Word: "dog", Translation: "perro", Category: "animals", CreatedAt: time.Now()},
		{Word: "apple", Translation: "manzana", Category: "food", CreatedAt: time.Now()},
	}

	if err := db.Create(&words).Error; err != nil {
		t.Fatalf("failed to seed words: %v", err)
	}

	// Create a user_word for the first word, set NextReview in the past so it's due today
	userWord := models.UserWord{
		WordID:            words[0].ID,
		BoxNumber:         1,
		LastReview:        time.Now().Add(-48 * time.Hour),
		NextReview:        time.Now().Add(-24 * time.Hour),
		CorrectAttempts:   0,
		IncorrectAttempts: 0,
	}

	if err := db.Create(&userWord).Error; err != nil {
		t.Fatalf("failed to seed user word: %v", err)
	}

	return words
}

func TestGetUserWordsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, db := setupTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	seedData(t, db)

	router := gin.New()
	router.GET("/userwords", handler.GetUserWords)

	req := httptest.NewRequest(http.MethodGet, "/userwords", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var got []models.UserWord
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(got) == 0 {
		t.Fatalf("expected at least one user word, got 0")
	}
}

func TestGetUserWordsByCategoryHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, db := setupTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	words := seedData(t, db)

	// Also add a user_word for the 'apple' word and make it due
	if err := db.Create(&models.UserWord{
		WordID:            words[2].ID,
		BoxNumber:         1,
		LastReview:        time.Now().Add(-48 * time.Hour),
		NextReview:        time.Now().Add(-24 * time.Hour),
		CorrectAttempts:   0,
		IncorrectAttempts: 0,
	}).Error; err != nil {
		t.Fatalf("failed to seed additional user word: %v", err)
	}

	router := gin.New()
	router.GET("/userwords/category/:category", handler.GetUserWordsByCategory)

	req := httptest.NewRequest(http.MethodGet, "/userwords/category/animals", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var got []models.UserWord
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// Expect only words that belong to the "animals" category
	if len(got) == 0 {
		t.Fatalf("expected user words for category 'animals', got none")
	}
	for _, uw := range got {
		if uw.Word.Category != "animals" {
			t.Fatalf("expected category 'animals', got %q for word id %d", uw.Word.Category, uw.WordID)
		}
	}
}

func TestUpdateUserWordHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, db := setupTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	words := seedData(t, db)

	// Ensure there's a user_word for words[0]
	var before models.UserWord
	if err := db.Where("word_id = ?", words[0].ID).First(&before).Error; err != nil {
		t.Fatalf("failed to fetch seeded user word: %v", err)
	}

	// Prepare request body to mark as learned (true)
	body := map[string]bool{"learned": true}
	bs, _ := json.Marshal(body)

	router := gin.New()
	router.PUT("/userwords/:wordID", handler.UpdateUserWord)

	req := httptest.NewRequest(http.MethodPut, "/userwords/"+strconv.FormatUint(uint64(words[0].ID), 10), bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	// Fetch the user_word again and check BoxNumber incremented (from 1 to 2)
	var after models.UserWord
	if err := db.Where("word_id = ?", words[0].ID).First(&after).Error; err != nil {
		t.Fatalf("failed to fetch user word after update: %v", err)
	}

	if after.BoxNumber <= before.BoxNumber {
		t.Fatalf("expected BoxNumber to increase after marking learned; before=%d after=%d", before.BoxNumber, after.BoxNumber)
	}
}

func TestSyncUserWordsAddsMissing(t *testing.T) {
	_, db := setupTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	// Seed words but do not create user_words for all of them
	words := []models.Word{
		{Word: "sun", Translation: "sol", Category: "nature", CreatedAt: time.Now()},
		{Word: "moon", Translation: "luna", Category: "nature", CreatedAt: time.Now()},
		{Word: "bread", Translation: "pan", Category: "food", CreatedAt: time.Now()},
	}
	if err := db.Create(&words).Error; err != nil {
		t.Fatalf("failed to seed words: %v", err)
	}

	// Only create a user_word for the first word
	if err := db.Create(&models.UserWord{
		WordID:            words[0].ID,
		BoxNumber:         1,
		LastReview:        time.Now(),
		NextReview:        time.Now(),
		CorrectAttempts:   0,
		IncorrectAttempts: 0,
	}).Error; err != nil {
		t.Fatalf("failed to seed initial user word: %v", err)
	}

	// Handler created by setupTest is not used here; recreate the handler wired to this DB.
	repo := repository.NewUserWordRepository(db)
	svc := services.NewUserWordService(repo)
	h := handlers.NewUserWordHandler(svc)

	// Call SyncUserWords which should add user_words for the missing words
	if err := h.SyncUserWords(); err != nil {
		t.Fatalf("SyncUserWords returned error: %v", err)
	}

	// Count user_words and compare to words length
	var count int64
	if err := db.Model(&models.UserWord{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to count user_words: %v", err)
	}

	if int(count) != len(words) {
		t.Fatalf("expected %d user_words after sync, got %d", len(words), count)
	}
}
