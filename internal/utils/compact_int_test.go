package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReserveAndWrite(t *testing.T) {
	t.Run("sufficient capacity", func(t *testing.T) {
		dst := make([]byte, 0, 10)
		err := reserveAndWrite(&dst, 3, func(pos int) {
			dst[pos+0] = 0xAA
			dst[pos+1] = 0xBB
			dst[pos+2] = 0xCC
		})
		require.NoError(t, err)
		require.Equal(t, 3, len(dst))
		want := []byte{0xAA, 0xBB, 0xCC}
		require.Equal(t, want, dst)
	})

	t.Run("insufficient capacity", func(t *testing.T) {
		dst := make([]byte, 0, 2)
		oldLen := len(dst)
		err := reserveAndWrite(&dst, 3, func(pos int) {
			dst[pos] = 0
		})
		require.Error(t, err)
		require.Equal(t, oldLen, len(dst))
	})
}

func TestAppendVarintUncompressed(t *testing.T) {
	t.Run("sufficient capacity int16", func(t *testing.T) {
		dst := make([]byte, 0, 2)
		err := appendVarintUncompressed(&dst, int16(0x1234))
		require.NoError(t, err)
		want := []byte{0x12, 0x34}
		require.Equal(t, want, dst)
	})

	t.Run("insufficient capacity int16", func(t *testing.T) {
		dst := make([]byte, 0, 1)
		oldLen := len(dst)
		err := appendVarintUncompressed(&dst, int16(0x1234))
		require.Error(t, err)
		require.Equal(t, oldLen, len(dst))
	})
}

func TestAppendVarintCompressed(t *testing.T) {
	t.Run("encode small int32", func(t *testing.T) {
		dst := make([]byte, 0, 4)
		err := appendVarintCompressed(&dst, int32(1))
		require.NoError(t, err)
		want := []byte{0x81, 0x02}
		require.Equal(t, want, dst)
	})

	t.Run("zero value is invalid", func(t *testing.T) {
		dst := make([]byte, 0, 2)
		err := appendVarintCompressed(&dst, int32(0))
		require.Error(t, err)
	})

	t.Run("insufficient capacity for compressed", func(t *testing.T) {
		dst := make([]byte, 0)
		oldLen := len(dst)
		err := appendVarintCompressed(&dst, int32(1))
		require.Error(t, err)
		require.Equal(t, oldLen, len(dst))
	})
}

func TestIsReadyForCompressedEncoding(t *testing.T) {
	require.False(t, isReadyForCompressedEncoding(int16(1)))

	// int32: small value -> true; value with top byte non-zero -> false
	require.True(t, isReadyForCompressedEncoding(int32(1)))
	require.False(t, isReadyForCompressedEncoding(int32(1<<24)))

	// int64: small value -> true; value with top byte non-zero -> false
	require.True(t, isReadyForCompressedEncoding(int64(1)))
	require.False(t, isReadyForCompressedEncoding(int64(1<<56)))
}

func TestRestoreVarIntFromBuffer_ErrorsAndRoundtrip(t *testing.T) {
	var dst64 int64

	t.Run("offset out of range", func(t *testing.T) {
		_, err := restoreVarIntFromBuffer([]byte{0x00}, -1, &dst64)
		require.ErrorIs(t, err, errOffsetOutOfRange)
	})

	t.Run("invalid prefix", func(t *testing.T) {
		buf := []byte{compactVarintMsb}
		_, err := restoreVarIntFromBuffer(buf, 0, &dst64)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid compact int prefix")
	})

	t.Run("insufficient payload", func(t *testing.T) {
		// prefix claims 3 data bytes but buffer is too small
		buf := []byte{byte(compactVarintPrefixBase + 3)}
		_, err := restoreVarIntFromBuffer(buf, 0, &dst64)
		require.ErrorIs(t, err, errInvalidCompactIntPayload)
	})

	t.Run("uncompressed restore", func(t *testing.T) {
		// big-endian bytes for 0x01020304
		buf := []byte{0x01, 0x02, 0x03, 0x04}
		var out int32
		next, err := restoreVarIntFromBuffer(buf, 0, &out)
		require.NoError(t, err)
		require.Equal(t, 4, next)
		require.Equal(t, int32(0x01020304), out)
	})

	t.Run("compressed restore roundtrip", func(t *testing.T) {
		dst := make([]byte, 0, 4)
		// produce compressed encoding for value 12345
		err := appendVarintCompressed(&dst, int32(12345))
		require.NoError(t, err)
		var out int32
		next, err := restoreVarIntFromBuffer(dst, 0, &out)
		require.NoError(t, err)
		require.Equal(t, int32(12345), out)
		require.Equal(t, len(dst), next)
	})
}

func TestAppendVarintCompressed_TooLarge(t *testing.T) {
	dst := make([]byte, 0)
	// very large value should result in an error mentioning 'compact varint too large'
	err := appendVarintCompressed(&dst, int64(math.MaxInt64))
	require.Error(t, err)
	require.Contains(t, err.Error(), "compact varint too large")
}
