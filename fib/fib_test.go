package fib

import (
	"testing"
)

var result int // private global variable to ensure benchmarks are actually run
type fibTest struct {
	value    int
	expected int
}

var fibTests = []fibTest{
	{1, 1},
	{2, 1},
	{3, 2},
	{4, 3},
	{5, 5},
	{6, 8},
	{7, 13},
	{8, 21},
	{9, 34},
}

func TestFibRecursive(t *testing.T) {
	for _, tt := range fibTests {
		actual := FibRecursive(tt.value)
		if actual != tt.expected {
			t.Errorf("FibRecursive(%d): expected %d, actual %d", tt.value, tt.expected, actual)
		}
	}
}

func TestFibIterative(t *testing.T) {
	for _, tt := range fibTests {
		actual := FibIterative(tt.value)
		if actual != tt.expected {
			t.Errorf("FibIterative(%d): expected %d, actual %d", tt.value, tt.expected, actual)
		}
	}
}

func TestFibRecursiveCache(t *testing.T) {
	for _, tt := range fibTests {
		actual := FibRecursiveCache(tt.value)
		if actual != tt.expected {
			t.Errorf("FibRecursiveCache(%d): expected %d, actual %d", tt.value, tt.expected, actual)
		}
	}
}

func TestFibTailRecursive(t *testing.T) {
	for _, tt := range fibTests {
		actual := FibTailRecursive(tt.value)
		if actual != tt.expected {
			t.Errorf("FibTailRecursive(%d): expected %d, actual %d", tt.value, tt.expected, actual)
		}
	}
}

func TestFibPowerMatrix(t *testing.T) {
	for _, tt := range fibTests {
		actual := FibPowerMatrix(tt.value)
		if actual != tt.expected {
			t.Errorf("FibPowerMatrix(%d): expected %d, actual %d", tt.value, tt.expected, actual)
		}
	}
}

func TestFibPowerMatrixRecursive(t *testing.T) {
	for _, tt := range fibTests {
		actual := FibPowerMatrixRecursive(tt.value)
		if actual != tt.expected {
			t.Errorf("FibPowerMatrixRecursive(%d): expected %d, actual %d", tt.value, tt.expected, actual)
		}
	}
}

func benchmarkFibRecursive(input int, b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = FibRecursive(input)
	}
	result = r
}

func BenchmarkFibRecursive1(b *testing.B)  { benchmarkFibRecursive(1, b) }
func BenchmarkFibRecursive2(b *testing.B)  { benchmarkFibRecursive(2, b) }
func BenchmarkFibRecursive4(b *testing.B)  { benchmarkFibRecursive(4, b) }
func BenchmarkFibRecursive8(b *testing.B)  { benchmarkFibRecursive(8, b) }
func BenchmarkFibRecursive16(b *testing.B) { benchmarkFibRecursive(16, b) }
func BenchmarkFibRecursive32(b *testing.B) { benchmarkFibRecursive(32, b) }

func benchmarkFibIterative(input int, b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = FibIterative(input)
	}
	result = r
}

func BenchmarkFibIterative1(b *testing.B)    { benchmarkFibIterative(1, b) }
func BenchmarkFibIterative2(b *testing.B)    { benchmarkFibIterative(2, b) }
func BenchmarkFibIterative4(b *testing.B)    { benchmarkFibIterative(4, b) }
func BenchmarkFibIterative8(b *testing.B)    { benchmarkFibIterative(8, b) }
func BenchmarkFibIterative16(b *testing.B)   { benchmarkFibIterative(16, b) }
func BenchmarkFibIterative32(b *testing.B)   { benchmarkFibIterative(32, b) }
func BenchmarkFibIterative64(b *testing.B)   { benchmarkFibIterative(64, b) }
func BenchmarkFibIterative128(b *testing.B)  { benchmarkFibIterative(128, b) }
func BenchmarkFibIterative256(b *testing.B)  { benchmarkFibIterative(256, b) }
func BenchmarkFibIterative512(b *testing.B)  { benchmarkFibIterative(512, b) }
func BenchmarkFibIterative1024(b *testing.B) { benchmarkFibIterative(1024, b) }

func benchmarkFibRecursiveCache(input int, b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = FibRecursiveCache(input)
	}
	result = r
}

func BenchmarkFibRecursiveCache1(b *testing.B)    { benchmarkFibRecursiveCache(1, b) }
func BenchmarkFibRecursiveCache2(b *testing.B)    { benchmarkFibRecursiveCache(2, b) }
func BenchmarkFibRecursiveCache4(b *testing.B)    { benchmarkFibRecursiveCache(4, b) }
func BenchmarkFibRecursiveCache8(b *testing.B)    { benchmarkFibRecursiveCache(8, b) }
func BenchmarkFibRecursiveCache16(b *testing.B)   { benchmarkFibRecursiveCache(16, b) }
func BenchmarkFibRecursiveCache32(b *testing.B)   { benchmarkFibRecursiveCache(32, b) }
func BenchmarkFibRecursiveCache64(b *testing.B)   { benchmarkFibRecursiveCache(64, b) }
func BenchmarkFibRecursiveCache128(b *testing.B)  { benchmarkFibRecursiveCache(128, b) }
func BenchmarkFibRecursiveCache256(b *testing.B)  { benchmarkFibRecursiveCache(256, b) }
func BenchmarkFibRecursiveCache512(b *testing.B)  { benchmarkFibRecursiveCache(512, b) }
func BenchmarkFibRecursiveCache1024(b *testing.B) { benchmarkFibRecursiveCache(1024, b) }

func benchmarkFibTailRecursive(input int, b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = FibTailRecursive(input)
	}
	result = r
}

func BenchmarkFibTailRecursive1(b *testing.B)    { benchmarkFibTailRecursive(1, b) }
func BenchmarkFibTailRecursive2(b *testing.B)    { benchmarkFibTailRecursive(2, b) }
func BenchmarkFibTailRecursive4(b *testing.B)    { benchmarkFibTailRecursive(4, b) }
func BenchmarkFibTailRecursive8(b *testing.B)    { benchmarkFibTailRecursive(8, b) }
func BenchmarkFibTailRecursive16(b *testing.B)   { benchmarkFibTailRecursive(16, b) }
func BenchmarkFibTailRecursive32(b *testing.B)   { benchmarkFibTailRecursive(32, b) }
func BenchmarkFibTailRecursive64(b *testing.B)   { benchmarkFibTailRecursive(64, b) }
func BenchmarkFibTailRecursive128(b *testing.B)  { benchmarkFibTailRecursive(128, b) }
func BenchmarkFibTailRecursive256(b *testing.B)  { benchmarkFibTailRecursive(256, b) }
func BenchmarkFibTailRecursive512(b *testing.B)  { benchmarkFibTailRecursive(512, b) }
func BenchmarkFibTailRecursive1024(b *testing.B) { benchmarkFibTailRecursive(1024, b) }

func benchmarkFibPowerMatrix(input int, b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = FibPowerMatrix(input)
	}
	result = r
}

func BenchmarkFibPowerMatrix1(b *testing.B)    { benchmarkFibPowerMatrix(1, b) }
func BenchmarkFibPowerMatrix2(b *testing.B)    { benchmarkFibPowerMatrix(2, b) }
func BenchmarkFibPowerMatrix4(b *testing.B)    { benchmarkFibPowerMatrix(4, b) }
func BenchmarkFibPowerMatrix8(b *testing.B)    { benchmarkFibPowerMatrix(8, b) }
func BenchmarkFibPowerMatrix16(b *testing.B)   { benchmarkFibPowerMatrix(16, b) }
func BenchmarkFibPowerMatrix32(b *testing.B)   { benchmarkFibPowerMatrix(32, b) }
func BenchmarkFibPowerMatrix64(b *testing.B)   { benchmarkFibPowerMatrix(64, b) }
func BenchmarkFibPowerMatrix128(b *testing.B)  { benchmarkFibPowerMatrix(128, b) }
func BenchmarkFibPowerMatrix256(b *testing.B)  { benchmarkFibPowerMatrix(256, b) }
func BenchmarkFibPowerMatrix512(b *testing.B)  { benchmarkFibPowerMatrix(512, b) }
func BenchmarkFibPowerMatrix1024(b *testing.B) { benchmarkFibPowerMatrix(1024, b) }

func benchmarkFibPowerMatrixRecursive(input int, b *testing.B) {
	var r int
	for i := 0; i < b.N; i++ {
		r = FibPowerMatrixRecursive(input)
	}
	result = r
}

func BenchmarkFibPowerMatrixRecursive1(b *testing.B)    { benchmarkFibPowerMatrixRecursive(1, b) }
func BenchmarkFibPowerMatrixRecursive2(b *testing.B)    { benchmarkFibPowerMatrixRecursive(2, b) }
func BenchmarkFibPowerMatrixRecursive4(b *testing.B)    { benchmarkFibPowerMatrixRecursive(4, b) }
func BenchmarkFibPowerMatrixRecursive8(b *testing.B)    { benchmarkFibPowerMatrixRecursive(8, b) }
func BenchmarkFibPowerMatrixRecursive16(b *testing.B)   { benchmarkFibPowerMatrixRecursive(16, b) }
func BenchmarkFibPowerMatrixRecursive32(b *testing.B)   { benchmarkFibPowerMatrixRecursive(32, b) }
func BenchmarkFibPowerMatrixRecursive64(b *testing.B)   { benchmarkFibPowerMatrixRecursive(64, b) }
func BenchmarkFibPowerMatrixRecursive128(b *testing.B)  { benchmarkFibPowerMatrixRecursive(128, b) }
func BenchmarkFibPowerMatrixRecursive256(b *testing.B)  { benchmarkFibPowerMatrixRecursive(256, b) }
func BenchmarkFibPowerMatrixRecursive512(b *testing.B)  { benchmarkFibPowerMatrixRecursive(512, b) }
func BenchmarkFibPowerMatrixRecursive1024(b *testing.B) { benchmarkFibPowerMatrixRecursive(1024, b) }
