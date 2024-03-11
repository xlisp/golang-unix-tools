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

```
```

