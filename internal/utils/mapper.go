package utils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Int64ToString formats an int64 to base-10 string without extra allocations.
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// StringToInt64 parses a base-10 int64 from string.
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ToNullString assumes input is already clean (no trimming) and converts
// to sql.NullString with minimal operations.
func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// FromNullString returns the underlying string or empty string when NULL.
func FromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// ToNullTime converts a *time.Time to sql.NullTime preserving UTC semantics.
func ToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	tt := t.UTC()
	return sql.NullTime{Time: tt, Valid: true}
}

// FromNullTime returns a pointer to UTC time or nil when invalid.
func FromNullTime(nt any) *time.Time {
	switch t := nt.(type) {
	case sql.NullTime:
		if t.Valid {
			tt := t.Time.UTC()
			return &tt
		}
	case gorm.DeletedAt:
		if t.Valid {
			tt := t.Time.UTC()
			return &tt
		}
	}
	return nil
}

func FromNullTimeRFC3339(nt any) string {
	tt := FromNullTime(nt)
	if tt != nil {
		return tt.Format(time.RFC3339)
	}
	return ""
}

func ToNullInt16(i int16) sql.NullInt16 {
	if i == 0 {
		return sql.NullInt16{Valid: false}
	}
	return sql.NullInt16{Int16: i, Valid: true}
}

func FromNullInt16(ni sql.NullInt16) int16 {
	if ni.Valid {
		return ni.Int16
	}
	return 0
}

// ToNullInt64 converts an int to sql.NullInt64.
func ToNullInt64(i int) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(i), Valid: true}
}

// FromNullInt64 converts sql.NullInt64 to int (zero when NULL).
func FromNullInt64(ni sql.NullInt64) int {
	if ni.Valid {
		return int(ni.Int64)
	}
	return 0
}

// ToNullBool converts bool to sql.NullBool (always valid).
func ToNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

// FromNullBool returns underlying bool or false when NULL.
func FromNullBool(nb sql.NullBool) bool {
	if nb.Valid {
		return nb.Bool
	}
	return false
}

// ToNullFloat converts float64 to sql.NullFloat64 (always valid).
func ToNullFloat(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: true}
}

// FromNullFloat returns underlying float64 or 0 when NULL.
func FromNullFloat(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0
}

// ToNullJSON marshals v to JSON and stores it in sql.NullString.
func ToNullJSON(v any) sql.NullString {
	if v == nil {
		return sql.NullString{Valid: false}
	}
	b, err := json.Marshal(v)
	if err != nil || len(b) == 0 {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: string(b), Valid: true}
}

// ToJSONBytes marshals v to JSON and returns raw bytes.
// When v is nil or marshals to "null", returns nil, nil.
func ToJSONBytes(v any) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 || bytes.Equal(b, nullJSONBytes) {
		return nil, nil
	}
	return b, nil
}

// FromJSONBytes unmarshals JSON bytes into dst (ptr). No-op on empty input.
func FromJSONBytes(b []byte, dst any) error {
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, dst)
}

// FromNullJSON unmarshals json string into dst (ptr) when valid.
func FromNullJSON(ns sql.NullString, dst any) error {
	if !ns.Valid {
		return nil
	}
	return json.Unmarshal([]byte(ns.String), dst)
}
