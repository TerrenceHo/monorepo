# Fibonacci Benchmarking

Here we go through several ways to implement the fibonacci sequence, including:
* Recursive
* Recursive with caching
* Tail Recursive
* Iterative
* Matrix Multiplcation
* \\(Log_n\\) Matrix Implementation.

I wrote this corresponding [blog
post](https://terrenceho.org/post/benchmarking-fibonacci/) to describe the methodology.

## Testing

Run `bazel test //fib:fib_test` to test fibonacci functions.

## Benchmarks

Run `bazel test //fib:fib_test --test_arg=-test.bench=. --test_output=all --test_timeout=100` 
to benchmark. You may need to up the timeout if your computer is especially
slow.
