package database_test

import (
	"testing"

	"github.com/brunojet/go-infra-backend/internal/database"
)

func TestNewSQLiteDatabase_DelegatesToAdapter(t *testing.T) {
	db, err := database.NewSQLiteDatabase()
	if err != nil {
		t.Fatalf("unexpected error creating sqlite database: %v", err)
	}
	if db == nil {
		t.Fatal("expected non-nil database")
	}
	if err := db.Close(); err != nil {
		t.Fatalf("unexpected error closing database: %v", err)
	}
}
