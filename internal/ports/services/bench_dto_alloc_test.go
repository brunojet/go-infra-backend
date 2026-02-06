package services

import (
	"testing"
)

var sample512 = func() []byte {
	s := make([]byte, 512)
	for i := range s {
		s[i] = byte(i % 256)
	}
	return s
}()

// Large value type: copying this copies the whole array (simulates worst-case DTO)
type LargeDTO struct {
	Data [512]byte
}

func makeLargeValue() LargeDTO {
	var d LargeDTO
	copy(d.Data[:], sample512)
	return d
}

func makeLargePointer() *LargeDTO {
	d := new(LargeDTO)
	copy(d.Data[:], sample512)
	return d
}

func fillLarge(out *LargeDTO) {
	copy(out.Data[:], sample512)
}

// String-based DTO: copying the struct copies the string header only (not the bytes)
type StringDTO struct {
	S string
}

func makeStringValue() StringDTO {
	return StringDTO{S: string(sample512)}
}

func makeStringPointer() *StringDTO {
	s := string(sample512)
	return &StringDTO{S: s}
}

func fillString(out *StringDTO) {
	out.S = string(sample512)
}

// Benchmarks for LargeDTO (value copies whole 512 bytes)
func Benchmark_ReturnValue_Large(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = makeLargeValue()
	}
}

func Benchmark_ReturnPointer_Large(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = makeLargePointer()
	}
}

func Benchmark_FillPointer_Large(b *testing.B) {
	b.ReportAllocs()
	var out LargeDTO
	for i := 0; i < b.N; i++ {
		fillLarge(&out)
	}
}

// Benchmarks for StringDTO (copying struct copies header only)
func Benchmark_ReturnValue_String(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = makeStringValue()
	}
}

func Benchmark_ReturnPointer_String(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = makeStringPointer()
	}
}

func Benchmark_FillPointer_String(b *testing.B) {
	b.ReportAllocs()
	var out StringDTO
	for i := 0; i < b.N; i++ {
		fillString(&out)
	}
}
