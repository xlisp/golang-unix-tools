package main

import (
	"fmt"
)

type PubSub struct {
	publishCh chan string
	subscribeCh chan chan string
}

func NewPubSub() *PubSub {
	return &PubSub{
		publishCh:   make(chan string),
		subscribeCh: make(chan chan string),
	}
}

func (ps *PubSub) Start() {
	subscribers := make([]chan string, 0)
	go func() {
		for {
			select {
			case msg := <-ps.publishCh:
				for _, sub := range subscribers {
					sub <- msg
				}
			case newSub := <-ps.subscribeCh:
				subscribers = append(subscribers, newSub)
			}
		}
	}()
}

func (ps *PubSub) Publish(msg string) {
	ps.publishCh <- msg
}

func (ps *PubSub) Subscribe() chan string {
	newSub := make(chan string)
	ps.subscribeCh <- newSub
	return newSub
}

func main() {
	ps := NewPubSub()
	ps.Start()

	sub1 := ps.Subscribe()
	sub2 := ps.Subscribe()

	ps.Publish("Hello, World!")

	fmt.Println(<-sub1)
	fmt.Println(<-sub2)
}