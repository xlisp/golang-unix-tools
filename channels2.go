package main
import "fmt"
import "time"
import "math/rand"

func boring(msg string) <-chan string { // Returns receive-only channel of strs.
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // Return channel to the caller
}

func main() {
	c:= boring("boring!") // Function returning a channel.
	for i := 0; i < 5; i++ {
		fmt.Printf("You say === : %q\n", <-c)
	}
	fmt.Println("You're boring; I'm leaving.")
}

// =>
// You say === : "boring! 0"
// You say === : "boring! 1"
// You say === : "boring! 2"
// You say === : "boring! 3"
// You say === : "boring! 4"
// You're boring; I'm leaving.
// 
