package utils

import (
	"encoding/csv"
	"fmt"
	"learning-cards/internal/models"
	"os"
	"time"
)

func ReadCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func ConvertToWords(records [][]string) ([]models.Word, error) {
	var words []models.Word

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) != 3 {
			return nil, fmt.Errorf("record on line %d: wrong number of fields, got %d, want 3", i+1, len(record))
		}

		word := models.Word{
			Word:        record[0],
			Translation: record[1],
			Category:    record[2],
			CreatedAt:   time.Now(),
		}

		words = append(words, word)
	}

	return words, nil
}
