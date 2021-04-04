// Fib returns the nth number in the Fibonacci series. This implements a series
// of Fib functions, tests, and benchmarks them. All implementations are
// correct, barring overflow issues, only differing in execution time.
package fib

// Basic recursive algorithm, with exponential run time.  Simply makes two extra
// calls for level. In testing benchmarks, we do not go above 32 for this
// function due to the extreme length of time it takes to complete, it is not
// worth benchmarking.  Needless to say, this function is terrible and slow.
func FibRecursive(n int) int {
	if n < 2 {
		return n
	}
	return FibRecursive(n-1) + FibRecursive(n-2)
}

// Optimized recursive option, utilizing a cache to eliminate repeated
// calculations.  This algorithm runs in linear time. Makes use of a helper
// function that recurses down to n = 1, and then builds the cache moving
// backwards. The number at the greatest index is the resulting answer.
func FibRecursiveCache(n int) int {
	cache := make([]int, n+1, n+1)
	fibRecursiveCache(n, &cache)
	return cache[n]
}

func fibRecursiveCache(n int, cache *[]int) {
	if n < 2 {
		(*cache)[0] = 0
		(*cache)[1] = 1
		return
	}
	fibRecursiveCache(n-1, cache)

	(*cache)[n] = (*cache)[n-1] + (*cache)[n-2]
}

// Another linear recursive implementation, without utilizing a memory cache.
// Instead each recursive call completes a portion of the calculation instead.
// Utilizes a recursive helper function to increment necessary values.
func FibTailRecursive(n int) int {
	return fibTailRecursive(n, 0, 1)
}

func fibTailRecursive(n, first, second int) int {
	if n == 0 {
		return first
	}
	return fibTailRecursive(n-1, second, first+second)
}

// Linear, iterative implementation.  Uses a for loop, and pre delcares temp
// variables to avoid initialization every loop.
func FibIterative(n int) int {
	var temp int
	first := 0
	second := 1
	for i := 0; i < n-1; i++ {
		temp = second
		second = first + second
		first = temp
	}
	return second
}

// Linear time execution, uses matrix multiplication to simplify logic
func FibPowerMatrix(n int) int {
	F := [2][2]int{
		[2]int{1, 1},
		[2]int{1, 0},
	}
	if n == 0 {
		return 0
	}
	fibPower(&F, n-1)
	return F[0][0]
}

func fibPower(F *[2][2]int, n int) {
	M := [2][2]int{
		[2]int{1, 1},
		[2]int{1, 0},
	}
	for i := 2; i <= n; i++ {
		fibMultiply(F, &M)
	}
}

func fibMultiply(F *[2][2]int, M *[2][2]int) {
	f := *F
	m := *M
	x := f[0][0]*m[0][0] + f[0][1]*m[1][0]
	y := f[0][0]*m[0][1] + f[0][1]*m[1][1]
	z := f[1][0]*m[0][0] + f[1][1]*m[1][0]
	w := f[1][0]*m[0][1] + f[1][1]*m[1][1]

	(*F)[0][0] = x
	(*F)[0][1] = y
	(*F)[1][0] = z
	(*F)[1][1] = w
}

func FibPowerMatrixRecursive(n int) int {
	F := [2][2]int{
		[2]int{1, 1},
		[2]int{1, 0},
	}

	if n == 0 {
		return 0
	}
	fibPowerRecursive(&F, n-1)
	return F[0][0]
}

func fibPowerRecursive(F *[2][2]int, n int) {
	if n == 0 || n == 1 {
		return
	}

	M := [2][2]int{
		[2]int{1, 1},
		[2]int{1, 0},
	}
	fibPowerRecursive(F, n/2)
	fibMultiply(F, F)
	if n%2 != 0 {
		fibMultiply(F, &M)
	}
}
