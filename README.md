# jim-emacs-fun-go

- [jim-emacs-fun-go](#jim-emacs-fun-go)
  - [gin](#gin)
  - [clojure core.async VS go goroutine](#clojure-coreasync-vs-go-goroutine)
  - [map](#map)
  - [函数作为参数](#%E5%87%BD%E6%95%B0%E4%BD%9C%E4%B8%BA%E5%8F%82%E6%95%B0)
  - [递归](#%E9%80%92%E5%BD%92)
  - [functor](#functor)
  - [调用C语言](#%E8%B0%83%E7%94%A8c%E8%AF%AD%E8%A8%80)
  - [run test](#run-test)
  - [add lib](#add-lib)

## gin

https://github.com/gin-gonic/gin

```
go mod init ginexample
go mod tidy

# export http_proxy=http://127.0.0.1:1087;export https_proxy=http://127.0.0.1:1087;

go get -u github.com/gin-gonic/gin

```

ginhttp.go
```go
package main

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
```

```
go run ginhttp.go
curl http://127.0.0.1:8080/ping
```
## clojure core.async VS go goroutine

* channels

```go
func boring(msg string, c chan string) {
	for i := 0; ; i++ {
		c <- fmt.Sprintf("%s %d", msg, i) // Expression to be sent can be any val.
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func main() {
	c := make(chan string)
	go boring("boring!", c)
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c) // Receive expression is just a value.
	}
	fmt.Println("You're boring; I'm leaving")
}
```

vs: clojure的>!! 和 go的c <- , clojure的<!!和 <-c

```clojure
(defn boring [msg c]
  (loop [i 0]
    (>!! c (str msg " " i))
    (recur (inc i))))

(defn -main [& args]
  (let [c (chan)]
    (go (boring "boring!" c))
    (dotimes [_ 5]
      (println (<!! c)))
    (println "You're boring; I'm leaving.")))
```

* function that returns a channel

go 宏把函数包起来

```go
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
		fmt.Printf("You say: %q\n", <-c)
	}
	fmt.Println("You're boring; I'm leaving.")
}
```

```clojure
(defn boring [msg]
  (let [c (chan)]
    (go (loop [i 0]
          (>!! c (str msg " " i))
          (Thread/sleep (rand-int 1000))
          (recur (inc i))))
    c))

(defn -main [& args]
  (let [c (boring "boring!")]
    (dotimes [_ 5]
      (println (<!! c)))
    (println "You're boring; I'm leaving.")))
```

## map

https://hedzr.com/golang/fp/golang-functional-programming-in-brief/

```go
// The higher-order-function takes an array and a function as arguments
func mapForEach(arr []string, fn func(it string) int) []int {
	var newArray = []int{}
	for _, it := range arr {
		// We are executing the method passed
		newArray = append(newArray, fn(it))
	}
	return newArray
}

func main() {
	var list = []string{"Orange", "Apple", "Banana", "Grape"}
	// we are passing the array and a function as arguments to mapForEach method.
	var out = mapForEach(list, func(it string) int {
		return len(it)
	})
	fmt.Println(out) // [6, 5, 6, 5]
}

```

## 函数作为参数

```go

package main

type Handler func (a int)

func xc(pa int, handler Handler) {
  handler(pa)
}

func main(){
  xc(123, func(a int){
	  print (a) //=> 123
  })
}

```
## 递归

```go
func factorialTailRecursive(num int) int {
	return factorial2(1, num)
}

func factorial2(accumulator, val int) int {
	if val == 1 {
		return accumulator
	}
	return factorial2(accumulator*val, val-1)
}

func main() {
	fmt.Println(factorialTailRecursive(10)) // 3628800
}

```

## functor

```go

package main

func add(a, b int) int { return a+b }
func sub(a, b int) int { return a-b }

var operators map[string]func(a, b int) int

func init(){
	operators = map[string]func(a, b int) int {
		"+": add,
			"-": sub,
	}
}


func calculator(a, b int, op string) int {
	fn := operators[op]
	return fn(a, b)
}

func main() {
	print(calculator(1, 2, "+")) //=> 3
}

```

## 调用C语言

```c
// add.c
#include <stdint.h>

void Add(uint64_t a, uint64_t b, uint64_t *ret) {
    *ret = a + b;
}
```
* call add.c
```go
package main

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
```


## run test
* `*_test.go`

```
go test
```

## add lib

```
$ go get github.com/rabbitmq/amqp091-go

go: downloading github.com/rabbitmq/amqp091-go v1.10.0
go: added github.com/rabbitmq/amqp091-go v1.10.0
```

