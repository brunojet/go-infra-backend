package utils

import (
	"database/sql"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertInt64ToT_Table(t *testing.T) {
	cases := []struct {
		name      string
		src       int64
		wantErr   bool
		wantInt32 int32
		wantInt16 int16
	}{
		{name: "int64 fits", src: int64(1234567890), wantErr: false},
		{name: "int32 max", src: int64(math.MaxInt32), wantErr: false},
		{name: "int32 overflow", src: int64(math.MaxInt32) + 1, wantErr: true},
		{name: "int16 fits", src: int64(32000), wantErr: false},
		{name: "int16 overflow", src: int64(math.MaxInt16) + 1, wantErr: true},
	}

	// test int64 passthrough
	got, err := convertInt64ToT[int64](cases[0].src)
	require.NoError(t, err)
	require.Equal(t, cases[0].src, got)

	// table-driven checks for int32 and int16 behavior
	for _, tc := range cases[1:] {
		t.Run(tc.name+"/int32", func(t *testing.T) {
			_, err := convertInt64ToT[int32](tc.src)
			switch tc.name {
			case "int32 overflow":
				require.Error(t, err)
			case "int16 overflow", "int16 fits":
				// these cases are for int16; skip
			default:
				require.NoError(t, err)
			}
		})

		t.Run(tc.name+"/int16", func(t *testing.T) {
			_, err := convertInt64ToT[int16](tc.src)
			switch tc.name {
			case "int16 overflow":
				require.Error(t, err)
			case "int32 overflow", "int32 max":
				// skip int32-specific
			case "int16 fits":
				require.NoError(t, err)
			}
		})
	}
}

func TestGetIntTypeLen_Table(t *testing.T) {
	cases := []struct {
		name string
		want int
	}{
		{name: "int16", want: int16Size},
		{name: "int32", want: int32Size},
		{name: "int64", want: int64Size},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.name {
			case "int16":
				require.Equal(t, tc.want, getIntTypeLen[int16]())
			case "int32":
				require.Equal(t, tc.want, getIntTypeLen[int32]())
			case "int64":
				require.Equal(t, tc.want, getIntTypeLen[int64]())
			}
		})
	}
}

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
	n := ToNullInt64(0)
	a.False(n.Valid)
	a.Equal(int64(0), FromNullInt64(n))

	n2 := ToNullInt64(5)
	a.True(n2.Valid)
	a.Equal(int64(5), FromNullInt64(n2))

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
