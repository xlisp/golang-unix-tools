package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

const (
	numTasks        = 1000000
	numPrimeChecks  = 10000
	fibonacciNumber = 40
)

func main() {
	start := time.Now()

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		computePrimes()
	}()

	go func() {
		defer wg.Done()
		computeFibonacci()
	}()

	go func() {
		defer wg.Done()
		performFloatingPointOperations()
	}()

	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("Total execution time: %v\n", duration)
	fmt.Printf("Number of CPUs: %d\n", numCPU)
	fmt.Printf("Performance score: %.2f (lower is better)\n", duration.Seconds())
}

func computePrimes() {
	count := 0
	for i := 2; i < numPrimeChecks; i++ {
		if isPrime(i) {
			count++
		}
	}
	fmt.Printf("Found %d prime numbers\n", count)
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func computeFibonacci() {
	result := fibonacci(fibonacciNumber)
	fmt.Printf("Fibonacci(%d) = %d\n", fibonacciNumber, result)
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func performFloatingPointOperations() {
	result := 0.0
	for i := 0; i < numTasks; i++ {
		result += math.Sin(float64(i)) * math.Cos(float64(i))
	}
	fmt.Printf("Floating-point operations result: %.2f\n", result)
}

/* Mac run :
 golang-unix-tools  master @ go build performance_test_speed.go

 golang-unix-tools  master @ ./performance_test_speed
Found 1229 prime numbers
Floating-point operations result: 0.20
Fibonacci(40) = 102334155
Total execution time: 577.379344ms
Number of CPUs: 8
Performance score: 0.58 (lower is better)

*/

/* Ubuntu PC run:

(base) ➜  golang-unix-tools git:(master) ./performance_test_speed
Found 1229 prime numbers
Floating-point operations result: 0.20
Fibonacci(40) = 102334155
Total execution time: 770.770276ms
Number of CPUs: 20
Performance score: 0.77 (lower is better)

*/

