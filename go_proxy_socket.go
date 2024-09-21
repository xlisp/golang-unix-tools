/*golang test how to use ` --socks5-hostname 127.0.0.1:1087` ? */
package main

import (
    "net/http"
    "testing"
    "time"
    "golang.org/x/net/proxy"
)

func TestWithSocks5Proxy(t *testing.T) {
    // Define the SOCKS5 proxy address
    socks5Proxy := "127.0.0.1:1087"

    // Create a SOCKS5 dialer
    dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
    if err != nil {
        t.Fatalf("Error creating SOCKS5 dialer: %v", err)
    }

    // Create a custom transport that uses the SOCKS5 dialer
    transport := &http.Transport{
        Dial: dialer.Dial,
    }

    // Create an HTTP client using the custom transport
    client := &http.Client{
        Transport: transport,
        Timeout:   10 * time.Second,
    }

    // Make a request through the SOCKS5 proxy
    resp, err := client.Get("http://www.google.com")
    if err != nil {
        t.Fatalf("Failed to connect through proxy: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status 200 OK, got %v", resp.Status)
    }

    t.Log("Request through SOCKS5 proxy was successful")
}

