// Removed the main function redeclaration and undefined variable issue.
package main

import (
	"testing"
)

func TestChannel(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 1
	value := <-ch
	if value != 1 {
		t.Errorf("Expected 1, but got %d", value)
	}
}