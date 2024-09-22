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
    serverAddr = "xxxxxxx:22"
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
