package utils

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type errMarshaler struct{}

func (e errMarshaler) MarshalJSON() ([]byte, error) { return nil, errors.New("marshal error") }

type sample struct{ A string }

func TestInt64String_Roundtrip(t *testing.T) {
	a := assert.New(t)
	s := Int64ToString(12345)
	v, err := StringToInt64(s)
	a.NoError(err)
	a.Equal(int64(12345), v)
}

func TestNullStringHelpers(t *testing.T) {
	a := assert.New(t)
	ns := ToNullString("")
	a.False(ns.Valid)
	a.Equal("", FromNullString(ns))

	ns2 := ToNullString("ok")
	a.True(ns2.Valid)
	a.Equal("ok", FromNullString(ns2))
}

func TestNullTimeHelpers(t *testing.T) {
	a := assert.New(t)
	var nilPtr *time.Time
	nt := ToNullTime(nilPtr)
	a.False(nt.Valid)
	a.Nil(FromNullTime(nt))

	now := time.Now()
	nt2 := ToNullTime(&now)
	a.True(nt2.Valid)
	got := FromNullTime(nt2)
	a.NotNil(got)
	a.Equal(now.UTC().Unix(), got.UTC().Unix())
}

func TestNullIntBoolFloat(t *testing.T) {
	a := assert.New(t)
	n := ToNullInt(0)
	a.False(n.Valid)
	a.Equal(0, FromNullInt(n))

	n2 := ToNullInt(5)
	a.True(n2.Valid)
	a.Equal(5, FromNullInt(n2))

	nb := ToNullBool(true)
	a.True(nb.Valid)
	a.True(FromNullBool(nb))

	nf := ToNullFloat(1.23)
	a.True(nf.Valid)
	a.Equal(1.23, FromNullFloat(nf))
}

func TestFromNullBool_Branches(t *testing.T) {
	// valid true
	nb := ToNullBool(true)
	assert.True(t, FromNullBool(nb))
	// invalid
	var nb2 sql.NullBool
	nb2.Valid = false
	assert.False(t, FromNullBool(nb2))
}

func TestFromNullFloat_Branches(t *testing.T) {
	nf := ToNullFloat(3.14)
	assert.Equal(t, 3.14, FromNullFloat(nf))

	var nf2 sql.NullFloat64
	nf2.Valid = false
	assert.Equal(t, 0.0, FromNullFloat(nf2))
}

func TestToNullJSON_ErrorsAndNil(t *testing.T) {
	// nil input -> invalid
	nn := ToNullJSON(nil)
	assert.False(t, nn.Valid)

	// marshaler error -> invalid
	nn = ToNullJSON(errMarshaler{})
	assert.False(t, nn.Valid)

	// valid struct -> valid
	n := ToNullJSON(sample{A: "x"})
	assert.True(t, n.Valid)
	assert.Contains(t, n.String, "x")
}

func TestToJSONBytes_Topaths(t *testing.T) {
	// nil input
	b, err := ToJSONBytes(nil)
	assert.NoError(t, err)
	assert.Nil(t, b)

	// pointer nil marshals to null -> returns nil
	var p *int = nil
	b, err = ToJSONBytes(p)
	assert.NoError(t, err)
	assert.Nil(t, b)

	// marshaler error
	b, err = ToJSONBytes(errMarshaler{})
	assert.Error(t, err)

	// valid
	b, err = ToJSONBytes(sample{A: "y"})
	assert.NoError(t, err)
	assert.NotNil(t, b)
}

func TestFromJSONBytes_Branches(t *testing.T) {
	var dst sample
	// empty input -> no-op
	assert.NoError(t, FromJSONBytes([]byte{}, &dst))

	// invalid json -> error
	err := FromJSONBytes([]byte("{invalid}"), &dst)
	assert.Error(t, err)

	// valid
	b, _ := ToJSONBytes(sample{A: "z"})
	assert.NoError(t, FromJSONBytes(b, &dst))
	assert.Equal(t, "z", dst.A)
}

func TestFromNullJSON_Branches(t *testing.T) {
	var dst sample
	// invalid nullstring -> no-op
	var ns sql.NullString
	ns.Valid = false
	assert.NoError(t, FromNullJSON(ns, &dst))

	// invalid json inside -> error
	ns = sql.NullString{String: "{bad}", Valid: true}
	assert.Error(t, FromNullJSON(ns, &dst))

	// valid
	ns = ToNullJSON(sample{A: "ok"})
	assert.NoError(t, FromNullJSON(ns, &dst))
	assert.Equal(t, "ok", dst.A)
}
