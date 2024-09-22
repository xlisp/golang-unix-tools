package main

import (
	"io"
    "io/ioutil"
    "log"
    "net"
    "time"

    "golang.org/x/crypto/ssh"
)

func main() {
    // SSH server details
    serverAddr := "xxxxxx:22"
    user := "ubuntu"
    keyPath := "/home/xlisp/xxxx.pem"

    // Read the private key file
    key, err := ioutil.ReadFile(keyPath)
    if err != nil {
        log.Fatalf("Unable to read private key: %v", err)
    }

    // Create the Signer for this private key
    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        log.Fatalf("Unable to parse private key: %v", err)
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
        log.Fatal("Failed to dial: ", err)
    }
    defer client.Close()
    log.Println("Connected successfully")

    // Set up remote forwarding
    log.Println("Setting up remote forwarding...")
    listener, err := client.Listen("tcp", "0.0.0.0:8899")
    if err != nil {
        log.Fatal("Failed to set up remote forwarding: ", err)
    }
    defer listener.Close()

    log.Println("Remote forwarding established. Listening on 0.0.0.0:8899")

    // Handle incoming connections
    for {
        log.Println("Waiting for incoming connection...")
        remoteConn, err := listener.Accept()
        if err != nil {
            log.Println("Failed to accept incoming connection:", err)
            continue
        }
        log.Printf("Accepted connection from %s\n", remoteConn.RemoteAddr())

        go handleConnection(remoteConn)
    }
}

func handleConnection(remoteConn net.Conn) {
    defer remoteConn.Close()

    // Connect to local port 22
    log.Println("Connecting to local SSH server...")
    localConn, err := net.Dial("tcp", "127.0.0.1:22")
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

/* run ssh  reverse  server ok: 
(base) xlisp@xlisp:~$ go run go_ssh_rev3.go
2024/09/22 12:31:34 Connecting to 54.222.239.188:22...
2024/09/22 12:31:35 Connected successfully
2024/09/22 12:31:35 Setting up remote forwarding...
2024/09/22 12:31:35 Remote forwarding established. Listening on 0.0.0.0:8899
2024/09/22 12:31:35 Waiting for incoming connection...
2024/09/22 12:31:50 Accepted connection from 127.0.0.1:42936
2024/09/22 12:31:50 Waiting for incoming connection...
2024/09/22 12:31:50 Connecting to local SSH server...
2024/09/22 12:31:50 Connected to local SSH server
2024/09/22 12:31:50 Starting bidirectional copy...

test ssh connect ok:
ubuntu@ip-172-31-2-171:~$ ssh -p 8899 xlisp@127.0.0.1
xlisp@127.0.0.1's password:
Welcome to Ubuntu 24.04.1 LTS (GNU/Linux 6.8.0-45-generic x86_64)
*/

