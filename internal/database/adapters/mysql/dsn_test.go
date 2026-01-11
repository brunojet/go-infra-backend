package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDSN_MinimalRequired(t *testing.T) {
	params := DSNParams{
		User:     "user",
		Password: "pass",
		Database: "db",
	}
	dsn, err := BuildDSN(params)
	assert.NoError(t, err)
	assert.Contains(t, dsn, "user:pass@tcp(localhost:3306)/db")
	assert.Contains(t, dsn, "charset=utf8mb4")
	assert.Contains(t, dsn, "collation=utf8mb4_unicode_ci")
}

func TestBuildDSN_AllFields(t *testing.T) {
	params := DSNParams{
		Host:     "myhost",
		Port:     1234,
		User:     "u",
		Password: "p",
		Database: "d",
		Options: map[string]string{
			"parseTime": "false",
			"loc":       "UTC",
			"custom":    "x",
		},
	}
	dsn, err := BuildDSN(params)
	assert.NoError(t, err)
	assert.Contains(t, dsn, "u:p@tcp(myhost:1234)/d")
	assert.Contains(t, dsn, "custom=x")
}

func TestBuildDSN_MissingUserOrDB(t *testing.T) {
	params := DSNParams{Password: "p"}
	_, err := BuildDSN(params)
	assert.Error(t, err)
	params = DSNParams{User: "u"}
	_, err = BuildDSN(params)
	assert.Error(t, err)
}

func TestBuildDSN_InvalidParseTime(t *testing.T) {
	params := DSNParams{
		User:     "u",
		Password: "p",
		Database: "d",
		Options:  map[string]string{"parseTime": "notabool"},
	}
	_, err := BuildDSN(params)
	assert.Error(t, err)
}

func TestBuildDSN_InvalidLoc(t *testing.T) {
	params := DSNParams{
		User:     "u",
		Password: "p",
		Database: "d",
		Options:  map[string]string{"loc": "notaloc"},
	}
	_, err := BuildDSN(params)
	assert.Error(t, err)
}

func TestBuildDSN_OptionsTrimming(t *testing.T) {
	params := DSNParams{
		User:     "u",
		Password: "p",
		Database: "d",
		Options:  map[string]string{"  custom  ": "  value  "},
	}
	dsn, err := BuildDSN(params)
	assert.NoError(t, err)
	assert.Contains(t, dsn, "custom=value")
}
