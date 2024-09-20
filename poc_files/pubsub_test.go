package main

import (
	"testing"
)

func TestPubSub(t *testing.T) {
	ps := NewPubSub()
	ps.Start()

	sub1 := ps.Subscribe()
	sub2 := ps.Subscribe()

	ps.Publish("Test Message")

	msg1 := <-sub1
	msg2 := <-sub2

	if msg1 != "Test Message" {
		t.Errorf("Expected \"Test Message\", but got %s", msg1)
	}

	if msg2 != "Test Message" {
		t.Errorf("Expected \"Test Message\", but got %s", msg2)
	}
}