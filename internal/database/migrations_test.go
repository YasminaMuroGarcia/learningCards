package database_test

import (
	"strings"
	"testing"

	dbpkg "learning-cards/internal/database"
	"learning-cards/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openInMemoryDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		// In some environments (CI or where CGO is disabled) the default sqlite driver
		// requires cgo and will return an error like:
		// "Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub"
		// Detect that case and skip the test instead of failing.
		if err != nil && (strings.Contains(err.Error(), "go-sqlite3 requires cgo") || strings.Contains(err.Error(), "CGO_ENABLED=0")) {
			t.Skipf("skipping sqlite-backed migration test due to environment: %v", err)
		}
		t.Fatalf("failed to open in-memory sqlite DB: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB from gorm DB: %v", err)
	}
	// ensure connection closed when test finishes
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}

func TestMigrateCreatesTablesAndColumns(t *testing.T) {
	db := openInMemoryDB(t)

	// Run the migrations from the package under test
	dbpkg.Migrate(db)

	// Verify tables were created
	if !db.Migrator().HasTable(&models.Word{}) {
		t.Fatalf("expected table for models.Word to exist after migration")
	}
	if !db.Migrator().HasTable(&models.UserWord{}) {
		t.Fatalf("expected table for models.UserWord to exist after migration")
	}

	// Verify some expected columns exist on the Word model.
	// Use struct field names as GORM checks them.
	expectedColumns := []string{"Word", "Translation", "Category", "CreatedAt"}
	for _, col := range expectedColumns {
		if !db.Migrator().HasColumn(&models.Word{}, col) {
			t.Fatalf("expected column %q on models.Word after migration", col)
		}
	}
}
