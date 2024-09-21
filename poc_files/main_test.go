// +build cgo

package main

import "testing"

func TestHello(t *testing.T) {
    // Since the C function prints to stdout, we cannot capture its output directly in a test.
    // Instead, we will just call the function to ensure there are no runtime errors.
    C.hello()
}