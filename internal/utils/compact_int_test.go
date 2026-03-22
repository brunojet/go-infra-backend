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
