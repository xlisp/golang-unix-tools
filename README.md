# jim-emacs-fun-go

## install

```
wget https://dl.google.com/go/go1.22.1.darwin-amd64.pkg
which go #=>  /usr/local/go/bin/go
```

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

channels

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
