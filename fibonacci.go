package main

import "fmt"

func fibonacci() func() int {
	a, b := 0, 1
	
	// 函数返回函数
	return func() int {
		a, b = b, a+b
		return a
	}
}

func main() {
	// 变量是函数的赋值写法
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
// 
