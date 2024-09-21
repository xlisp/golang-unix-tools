// main.go
package main

/*
#include "add.c"

// Declare the C function
extern void Add(uint64_t a, uint64_t b, uint64_t *ret);
*/
import "C"
import (
    "fmt"
)

func main() {
    var result C.uint64_t

    // Call the C function
    C.Add(10, 20, &result)

    fmt.Println("Result:", result)
}

/* Run: 
$ go run addgo.go
Result: 30
*/

