package main

import (
	"net"
	"testing"
	"time"
)

func TestServerStart(t *testing.T) {
	// Start the server in a separate goroutine
	go main()

	// Wait for the server to start
	time.Sleep(10 * time.Millisecond)

	// Attempt to connect to port 6379
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	conn.Close()

}