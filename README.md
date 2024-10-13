# jim-emacs-fun-go
* System underlying scripts (proxy, intranet penetration~), Docker hacker, cuda and other C/C++ library calls, Unix high-performance tools, etc.


- [jim-emacs-fun-go](#jim-emacs-fun-go)
  - [defer func](#defer-func)
  - [go struct](#go-struct)
  - [gin](#gin)
  - [clojure core.async VS go goroutine](#clojure-coreasync-vs-go-goroutine)
  - [map](#map)
  - [Function as Parameter](#function-as-parameter)
  - [Recursion](#recursion)
  - [functor](#functor)
  - [Calling C Language](#calling-c-language)
  - [run test](#run-test)
  - [add lib](#add-lib)
  - [AST Parsing Function Relationships](#ast-parsing-function-relationships)
  - [SSH Tunneling and Auto Reconnect](#ssh-tunneling-and-auto-reconnect)

## defer func

when function do finished, will do defer

```go
func (wrapper *Wrapper) buildSomthing() {

... or

	wrapper.Lock()
	defer wrapper.Unlock()

...
	wrapper.Lock()
	defer func() {
		wrapper.Unlock()
		...
	}()
...
}
```

## go struct 

```go
type Wrapper struct {
        mqToChannel <-chan string
        someStatus bool
}
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

* [more examples](./go_vs_clojure_async.md)

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

vs: clojure's >!! and go's c <- , clojure's <!! and <-c

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

// Function returns another function

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

* loop get token

```clj
(mount/defstate functor-official-account-token-updator
  :start (let [poison (a/promise-chan)]
           (a/go-loop [tout (a/timeout 0)]
             (let [[_ port] (a/alts! [poison tout] :priority true)]
               (when-not (= port poison)
                 (let [{:keys [appid secret]} (:functortech @config/functor-api-conf)
                       {:keys [access_token expires_in] :as res}
                       (:body (client/get (str
                                            "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="
                                            appid
                                            "&secret="
                                            secret)
                                {:accept :json :as :json}))]
                   (reset! functor-official-account-token access_token)
                   (if expires_in
                     (info "functor-official-account access token expires in" expires_in res)
                     (error "functor-official-account access token error" res))
                   (recur (a/timeout (* 1000 (if expires_in (/ expires_in 3) 600))))))))
           poison)
  :stop (a/close! @functor-official-account-token-updator))
```

1. Atomic Value: We use atomic.Value for thread-safe token updates.
2. Token Fetching: fetchToken() makes the HTTP request to the WeChat API and decodes the response.
3. Periodic Token Refresh: startTokenUpdater() runs a goroutine that periodically fetches the token. It calculates the sleep duration based on the expiration time (expires_in).
4. Graceful Stop: The updater can be stopped by closing the stop channel, mimicking the :stop logic from the Clojure version.

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	// Atomic string to hold the token
	functorOfficialAccountToken atomic.Value
	// Configuration holding the appid and secret
	config = struct {
		AppID  string
		Secret string
	}{
		AppID:  "your_appid_here",
		Secret: "your_secret_here",
	}
)

// Response structure from the WeChat token API
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// fetchToken retrieves the token from the WeChat API
func fetchToken() (TokenResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", config.AppID, config.Secret)
	resp, err := http.Get(url)
	if err != nil {
		return TokenResponse{}, err
	}
	defer resp.Body.Close()

	var result TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return TokenResponse{}, err
	}
	return result, nil
}

// startTokenUpdater starts a goroutine that fetches the token periodically
func startTokenUpdater() chan struct{} {
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				// Fetch the token
				tokenResp, err := fetchToken()
				if err != nil {
					log.Printf("Error fetching token: %v", err)
					time.Sleep(10 * time.Minute) // Retry after 10 minutes if error
					continue
				}

				// Store the token
				functorOfficialAccountToken.Store(tokenResp.AccessToken)
				log.Printf("Token fetched successfully, expires in %d seconds", tokenResp.ExpiresIn)

				// Calculate next update time as a third of the expiration time, or default to 10 minutes if not provided
				updateInterval := time.Duration(tokenResp.ExpiresIn/3) * time.Second
				if updateInterval <= 0 {
					updateInterval = 10 * time.Minute
				}

				time.Sleep(updateInterval)
			}
		}
	}()
	return stop
}

func stopTokenUpdater(stop chan struct{}) {
	close(stop)
}

func main() {
	// Start the token updater
	stop := startTokenUpdater()

	// Simulate running for a while
	time.Sleep(30 * time.Minute)

	// Stop the token updater
	stopTokenUpdater(stop)

	log.Println("Token updater stopped.")
}

```

## Function as Parameter

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
## Recursion

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

## Calling C Language

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
* run all path `*_test.go`
```
go test ./...
```
* run current path `./*_test.go`
```
cd path && go test
```

## add lib

```
$ go get github.com/rabbitmq/amqp091-go

go: downloading github.com/rabbitmq/amqp091-go v1.10.0
go: added github.com/rabbitmq/amqp091-go v1.10.0
```

## AST Parsing Function Relationships

![Detailed tutorial on viewing function call relationships in a Go project](./show_fun_refs_project.md)

```go
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

var functionCalls = make(map[string][]string)

// Parse the function calls in a function declaration
func inspectFuncDecl(node ast.Node) bool {
	// We are only interested in function declarations
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return true
	}

	funcName := funcDecl.Name.Name

	// Traverse the function body to find function calls
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Get the name of the called function
		switch fun := callExpr.Fun.(type) {
		case *ast.Ident:
			// Simple function call
			functionCalls[funcName] = append(functionCalls[funcName], fun.Name)
		case *ast.SelectorExpr:
			// Method call (e.g., obj.Method)
			functionCalls[funcName] = append(functionCalls[funcName], fun.Sel.Name)
		}
		return true
	})

	return true
}

// Generate the Graphviz DOT format output
func generateDot() {
	fmt.Println("digraph G {")
	for caller, callees := range functionCalls {
		for _, callee := range callees {
			fmt.Printf("    \"%s\" -> \"%s\";\n", caller, callee)
		}
	}
	fmt.Println("}")
}

// Parse all Go files in the given directory
func parseGoFilesInDir(dir string) {
	fs := token.NewFileSet()

	// Walk through the directory to find Go files
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .go files
		if filepath.Ext(path) == ".go" {
			node, err := parser.ParseFile(fs, path, nil, 0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
				return err
			}

			// Walk through the AST and inspect function declarations
			ast.Inspect(node, inspectFuncDecl)
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking through directory: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_directory>")
		return
	}
	// Parse the directory
	dir := os.Args[1]
	parseGoFilesInDir(dir)

	// Generate Graphviz DOT output
	generateDot()
}

	// =>> Refactored
```


## SSH Tunneling and Auto Reconnect

![Detailed usage tutorial](https://github.com/xlisp/golang-unix-tools/blob/master/go_ssh_reverse_proxy.md)

```go
package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net"
    "strings"
    "time"

    "golang.org/x/crypto/ssh"
)

const (
    serverAddr = "xxxxx:22"
    user       = "ubuntu"
    keyPath    = "/home/xlisp/xxx.pem"
    localPort  = "127.0.0.1:22"
    remotePort = "0.0.0.0:8899"
)

func main() {
    for {
        if err := runSSHForward(); err != nil {
            log.Printf("SSH forwarding stopped: %v", err)
            log.Println("Attempting to reconnect in 5 seconds...")
            time.Sleep(5 * time.Second)
        }
    }
}

func runSSHForward() error {
    // Read the private key file
    key, err := ioutil.ReadFile(keyPath)
    if err != nil {
        return fmt.Errorf("unable to read private key: %v", err)
    }

    // Create the Signer for this private key
    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        return fmt.Errorf("unable to parse private key: %v", err)
    }

    // Create SSH config
    config := &ssh.ClientConfig{
        User: user,
        Auth: []ssh.AuthMethod{
            ssh.PublicKeys(signer),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout:         15 * time.Second,
    }

    // Connect to SSH server
    log.Printf("Connecting to %s...\n", serverAddr)
    client, err := ssh.Dial("tcp", serverAddr, config)
    if err != nil {
        return fmt.Errorf("failed to dial: %v", err)
    }
    defer client.Close()
    log.Println("Connected successfully")

    // Set up remote forwarding
    log.Println("Setting up remote forwarding...")
    listener, err := client.Listen("tcp", remotePort)
    if err != nil {
        return fmt.Errorf("failed to set up remote forwarding: %v", err)
    }
    defer listener.Close()

    log.Printf("Remote forwarding established. Listening on %s\n", remotePort)

    // Start health check
    go healthCheck(client)

    // Handle incoming connections
    for {
        log.Println("Waiting for incoming connection...")
        remoteConn, err := listener.Accept()
        if err != nil {
            return fmt.Errorf("failed to accept incoming connection: %v", err)
        }
        log.Printf("Accepted connection from %s\n", remoteConn.RemoteAddr())

        go handleConnection(remoteConn)
    }
}

func handleConnection(remoteConn net.Conn) {
    defer remoteConn.Close()

    // Connect to local port 22
    log.Println("Connecting to local SSH server...")
    localConn, err := net.Dial("tcp", localPort)
    if err != nil {
        log.Println("Failed to connect to local port:", err)
        return
    }
    defer localConn.Close()
    log.Println("Connected to local SSH server")

    // Copy data between connections
    log.Println("Starting bidirectional copy...")
    go func() {
        _, err := io.Copy(remoteConn, localConn)
        if err != nil {
            log.Println("Error copying data from local to remote:", err)
        }
    }()

    _, err = io.Copy(localConn, remoteConn)
    if err != nil {
        log.Println("Error copying data from remote to local:", err)
    }
    log.Println("Connection closed")
}

func healthCheck(client *ssh.Client) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        // Run netstat command on remote server
        session, err := client.NewSession()
        if err != nil {
            log.Printf("Failed to create session: %v", err)
            continue
        }
        defer session.Close()

        output, err := session.CombinedOutput("netstat -anp | grep 8899")
        if err != nil {
            log.Printf("Failed to run netstat command: %v", err)
            continue
        }

        // Check if the port is being listened to
        if strings.Contains(string(output), "LISTEN") {
            log.Println("Health check passed: Port 8899 is being listened to")
        } else {
            log.Println("Health check failed: Port 8899 is not being listened to")
        }
    }
}
