package main

func add(a, b int) int { return a + b }
func sub(a, b int) int { return a - b }

var operators map[string]func(a, b int) int

func init() {
	operators = map[string]func(a, b int) int{
		"+": add,
		"-": sub,
	}
}

func calculator(a, b int, op string) int {
	//  ./functor.go:18:30: invalid operation: op && fn != nil (mismatched types string and untyped bool)
	//if fn, ok := operators[op]; op && fn!=nil{
	//	return fn(a, b)
	//}
	//return 0

	// =>> Refactored
	fn := operators[op]
	return fn(a, b)

}

func main() {

	print(calculator(1, 2, "+")) //=> 3

}
