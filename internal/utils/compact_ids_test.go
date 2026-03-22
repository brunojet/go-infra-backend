package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeCompositeKey_Table(t *testing.T) {
	t.Run("no values error", func(t *testing.T) {
		_, err := EncodeCompositeKey[int32]()
		require.ErrorIs(t, err, errCompositeKeyValuesRequired)
	})

	t.Run("invalid decode string", func(t *testing.T) {
		var out int32
		err := DecodeCompositeKey("!!!notbase64!!!", &out)
		require.Error(t, err)
	})

	t.Run("roundtrip encode/decode", func(t *testing.T) {
		a, b, c := int32(1), int32(2), int32(3000)
		s, err := EncodeCompositeKey(a, b, c)
		require.NoError(t, err)
		var ra, rb, rc int32
		err = DecodeCompositeKey(s, &ra, &rb, &rc)
		require.NoError(t, err)
		require.Equal(t, a, ra)
		require.Equal(t, b, rb)
		require.Equal(t, c, rc)
	})
}

func TestEncodeCompositeKey_Errors(t *testing.T) {
	// srcs vazio
	_, err := EncodeCompositeKey[int32]()
	if err == nil {
		t.Error("esperado erro para srcs vazio")
	}
}

func TestDecodeCompositeKey_Errors(t *testing.T) {
	// dsts vazio
	err := DecodeCompositeKey[int32]("abc")
	if err == nil {
		t.Error("esperado erro para dsts vazio")
	}

	// base64 inválido
	var v int32
	err = DecodeCompositeKey("!!!", &v)
	if err == nil {
		t.Error("esperado erro para base64 inválido")
	}
}
