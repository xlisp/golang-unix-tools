package main
import "fmt"
import "time"
import "math/rand"

// https://gist.github.com/xyproto/6584125#using-channels-1201

func boring(msg string, c chan string) {
	for i := 0; ; i++ {
		c <- fmt.Sprintf("%s %d", msg, i) // Expression to be sent can be any val.
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func main() {
	//now := time.Now()
	//fmt.Println("Now() :" + now.String()) //=> Now() :2024-03-11 17:48:36.757627 +0800 CST m=+0.000146860
	c := make(chan string)
	go boring("boring!", c)
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c) // Receive expression is just a value.
	}
	fmt.Println("You're boring; I'm leaving")
}
// =>
// You say: "boring! 0"
// You say: "boring! 1"
// You say: "boring! 2"
// You say: "boring! 3"
// You say: "boring! 4"
// You're boring; I'm leaving
// 
