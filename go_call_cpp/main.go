package main

/*
#cgo CXXFLAGS: -std=c++11
#cgo LDFLAGS: -lstdc++
#include <stdlib.h>

// Declare the C++ functions to be used
extern int Add(int a, int b);
extern void PrintMessage();
*/
import "C"
import "fmt"

func main() {
    // Call the Add function from C++
    result := C.Add(3, 5)
    fmt.Printf("Result of Add(3, 5): %d\n", result)

    // Call the PrintMessage function from C++
    C.PrintMessage()
}

// g++ -c -o cppcode.o cppcode.cpp
// go build -o main main.go
// main


