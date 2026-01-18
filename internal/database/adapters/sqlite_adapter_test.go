package adapters_test

import (
	"testing"

	dbadpt "github.com/brunojet/go-infra-backend/internal/database/adapters"
	dbcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"
)

func TestNewInMemory_ReturnsDatabase(t *testing.T) {
	db, err := dbadpt.NewInMemory()
	if err != nil {
		t.Fatalf("unexpected error creating in-memory sqlite: %v", err)
	}
	if db == nil {
		t.Fatal("expected non-nil database")
	}

	// ensure it satisfies the contract
	var _ dbcontracts.Database = db

	// ensure GormDB is available
	gdb := db.GormDB()
	if gdb == nil {
		t.Fatal("expected non-nil gorm DB")
	}

	// Close should succeed
	if err := db.Close(); err != nil {
		t.Fatalf("unexpected error closing db: %v", err)
	}
}
