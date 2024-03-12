package main

import "fmt"

func factorial(num int) int {
	result := 1
	for ; num > 0; num-- {
		result *= num
	}
	return result
}

//func main() {
//	fmt.Println(factorial(10)) // 3628800
//}

// => FP
func factorialTailRecursive(num int) int {
	return factorial2(1, num)
}

func factorial2(accumulator, val int) int {
	if val == 1 {
		return accumulator
	}
	return factorial2(accumulator*val, val-1)
}

func main() {
	fmt.Println(factorialTailRecursive(10)) // 3628800
}
