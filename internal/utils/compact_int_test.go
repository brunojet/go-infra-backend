package utils

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type intTestCase[T AnyInt] struct {
	name        string
	value       T
	valid       bool
	wantBufSize int
}

func TestAppendAndReadVarint_Int16(t *testing.T) {
	testCases := []intTestCase[int16]{
		{"min", 0, true, 1},
		{"max", math.MaxInt16, true, 2},
		{"mid", 12345, true, 2},
		{"neg1", -1, false, 0},
	}
	testAppendAndReadVarint(t, testCases, 3)
}

func TestAppendAndReadVarint_Int32(t *testing.T) {
	testCases := []intTestCase[int32]{
		{"min", 0, true, 1},
		{"max", math.MaxInt32, true, 4},
		{"mid", 123456789, true, 4},
		{"neg1", -1, false, 0},
	}
	testAppendAndReadVarint(t, testCases, 3)
}

func TestAppendAndReadVarint_Int64(t *testing.T) {
	testCases := []intTestCase[int64]{
		{"min", 0, true, 1},
		{"max", math.MaxInt64, true, 8},
		{"mid", 9876543210123, true, 7},
		{"neg1", -1, false, 0},
	}
	testAppendAndReadVarint(t, testCases, 3)
}

func testAppendAndReadVarint[T AnyInt](t *testing.T, testCases []intTestCase[T], loops int) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%v", tc.name, tc.value), func(t *testing.T) {
			needed := getIntTypeLen[T]() * loops
			buf := make([]byte, 0, needed)
			err := error(nil)

			for range loops {
				buf, err = appendVarintToBuffer(buf, tc.value)
				if !tc.valid {
					assert.Error(t, err)
					return
				} else {
					assert.NoError(t, err)
				}
			}

			for range loops {
				var consumed int
				var out T
				consumed, err = readVarintFromBufferAt(buf, consumed, &out)
				assert.NoError(t, err)
				assert.Equal(t, tc.value, out)
			}
		})
	}
}

func TestAppendVarint_BufferInsuficiente(t *testing.T) {
	buf := make([]byte, 0, 1) // capacidade insuficiente
	_, err := appendVarintToBuffer(buf, int32(12345))
	assert.Error(t, err)
}

func TestAppendAndReadVarint_SequenceDifferentValues(t *testing.T) {
	// Testa roundtrip de valores diferentes em sequência, igual ao uso em compact_ids
	values := []int32{1, 2, 3000, 0, 127, 128, 255, 256, 32767}
	buf := make([]byte, 0, len(values)*getIntTypeLen[int32]())
	for _, v := range values {
		var err error
		buf, err = appendVarintToBuffer(buf, v)
		assert.NoError(t, err)
	}
	offset := 0
	for i, want := range values {
		var got int32
		next, err := readVarintFromBufferAt(buf, offset, &got)
		assert.NoError(t, err, "erro ao ler valor %d", i)
		assert.Equal(t, want, got, "valor lido diferente do esperado no indice %d", i)
		offset = next
	}
	assert.Equal(t, len(buf), offset, "offset final deve consumir todo buffer")
}

func TestReadVarintFromBufferAt_Errors(t *testing.T) {
	// offset > bufLen
	var out int32
	buf := []byte{0x00}
	_, err := readVarintFromBufferAt(buf, 10, &out)
	if err == nil {
		t.Error("esperado erro de offset fora do range")
	}

	// buffer insuficiente para valor esperado
	buf = []byte{0x01, 0xFF} // header indica 1 byte, mas needed pode ser maior
	_, err = readVarintFromBufferAt(buf, 1, &out)
	if err == nil {
		t.Error("esperado erro de buffer insuficiente")
	}

	// header inválido: realBytes+1 > needed
	// Para int16, needed=2, então realBytes=2 (header=0x82) => 2+1 > 2
	var out16 int16
	buf = []byte{0x82, 0x01, 0x02}
	_, err = readVarintFromBufferAt(buf, 0, &out16)
	if err == nil {
		t.Error("esperado erro de varint compactado inválido")
	}
}

func BenchmarkAppendVarintToBuffer_Int32(b *testing.B) {
	buf := make([]byte, 0, 8)
	for i := 0; b.Loop(); i++ {
		_, err := appendVarintToBuffer(buf[:0], int32(i))
		if err != nil {
			b.Fatalf("erro: %v", err)
		}
	}
}

func BenchmarkReadVarintFromBufferAt_Int32(b *testing.B) {
	buf, _ := appendVarintToBuffer(make([]byte, 0, 8), int32(123456))
	var out int32
	for b.Loop() {
		_, err := readVarintFromBufferAt(buf, 0, &out)
		if err != nil {
			b.Fatalf("erro: %v", err)
		}
	}
}
