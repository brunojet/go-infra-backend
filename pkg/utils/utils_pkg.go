package utils

import (
	"database/sql"
	"time"

	"github.com/brunojet/go-infra-backend/internal/utils"
)

// Re-exported types
type AnyInt = utils.AnyInt

// Composite key helpers
func EncodeCompositeKey[T AnyInt](srcs ...T) (string, error) {
	return utils.EncodeCompositeKey(srcs...)
}

func DecodeCompositeKey[T AnyInt](encoded string, dsts ...*T) error {
	return utils.DecodeCompositeKey(encoded, dsts...)
}

// JSON helpers
func ToJSONBytes(v any) ([]byte, error)             { return utils.ToJSONBytes(v) }
func FromJSONBytes(b []byte, dst any) error         { return utils.FromJSONBytes(b, dst) }
func ToNullJSON(v any) sql.NullString               { return utils.ToNullJSON(v) }
func FromNullJSON(ns sql.NullString, dst any) error { return utils.FromNullJSON(ns, dst) }

// Simple mappers
func Int64ToString(i int64) string          { return utils.Int64ToString(i) }
func StringToInt64(s string) (int64, error) { return utils.StringToInt64(s) }

func ToNullString(s string) sql.NullString    { return utils.ToNullString(s) }
func FromNullString(ns sql.NullString) string { return utils.FromNullString(ns) }

// Time mappers
func ToNullTime(t *time.Time) sql.NullTime { return utils.ToNullTime(t) }
func FromNullTime(nt any) *time.Time       { return utils.FromNullTime(nt) }
func FromNullTimeRFC3339(nt any) string    { return utils.FromNullTimeRFC3339(nt) }

func ToNullInt16(i int16) sql.NullInt16    { return utils.ToNullInt16(i) }
func FromNullInt16(ni sql.NullInt16) int16 { return utils.FromNullInt16(ni) }

func ToNullInt64(i int) sql.NullInt64    { return utils.ToNullInt64(i) }
func FromNullInt64(ni sql.NullInt64) int { return utils.FromNullInt64(ni) }

func ToNullBool(b bool) sql.NullBool    { return utils.ToNullBool(b) }
func FromNullBool(nb sql.NullBool) bool { return utils.FromNullBool(nb) }

func ToNullFloat(f float64) sql.NullFloat64    { return utils.ToNullFloat(f) }
func FromNullFloat(nf sql.NullFloat64) float64 { return utils.FromNullFloat(nf) }
