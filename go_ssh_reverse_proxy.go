package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	// SSH server credentials
	server := "xxxxxxx:22"
	user := "ubuntu"

	// Path to private key file
	keyPath := "/path/to/private/key"

	// Read private key
	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// Establish an SSH connection
	client, err := ssh.Dial("tcp", server, config)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer client.Close()

	// Listen on local port 8899 for reverse tunnel
	localListener, err := net.Listen("tcp", "0.0.0.0:8899")
	if err != nil {
		log.Fatalf("failed to listen on local port 8899: %v", err)
	}
	defer localListener.Close()

	fmt.Println("Listening on port 8899 for reverse SSH tunnel...")

	for {
		// Accept a new connection from the listener
		localConn, err := localListener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		// Forward this connection through the SSH tunnel
		remoteConn, err := client.Dial("tcp", "0.0.0.0:22")
		if err != nil {
			log.Printf("failed to establish remote connection: %v", err)
			localConn.Close()
			continue
		}

		// Pipe data between local and remote connections
		go func() {
			defer localConn.Close()
			defer remoteConn.Close()
			err = pipe(localConn, remoteConn)
			if err != nil {
				log.Printf("error piping data: %v", err)
			}
		}()
	}
}

// Pipe forwards data between two network connections
func pipe(conn1, conn2 net.Conn) error {
	errChan := make(chan error, 2)

	copyConn := func(dst, src net.Conn) {
		_, err := net.Copy(dst, src)
		errChan <- err
	}

	go copyConn(conn1, conn2)
	go copyConn(conn2, conn1)

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}
