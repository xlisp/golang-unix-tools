package main

import "fmt"

func fibonacci() func() int {
	a, b := 0, 1

	// Function returns a function
	return func() int {
		a, b = b, a+b
		return a
	}
}

func main() {
	// Assigning a function to a variable
	f := fibonacci()

	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

// ==>>
// 1
// 1
// 2
// 3
// 5
// 8
// 13
// 21
// 34
// 55
