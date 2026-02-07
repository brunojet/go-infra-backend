package utils

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestJSONHelpers(t *testing.T) {
	a := assert.New(t)
	// ToJSONBytes nil
	b, err := ToJSONBytes(nil)
	a.NoError(err)
	a.Nil(b)

	// ToJSONBytes valid
	m := map[string]string{"x": "y"}
	jb, err := ToJSONBytes(m)
	a.NoError(err)
	a.NotNil(jb)

	var out map[string]string
	a.NoError(FromJSONBytes(jb, &out))
	a.Equal("y", out["x"])

	// FromNullJSON
	ns := sql.NullString{String: string(jb), Valid: true}
	var out2 map[string]string
	a.NoError(FromNullJSON(ns, &out2))
	a.Equal("y", out2["x"])

	// ToNullJSON nil
	nn := ToNullJSON(nil)
	a.False(nn.Valid)

	// ToNullJSON value
	nn2 := ToNullJSON(m)
	a.True(nn2.Valid)
}
