package utils

import (
	"math"
	"testing"

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
