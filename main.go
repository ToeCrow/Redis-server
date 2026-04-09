package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// step 1: bind to port 6379
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error at start:", err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Redis-clone listen to :6379")

	for {
		// step 2: wait for client connection (ex. redis-cli)
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// step 3: handle client connection in a separate goroutine (now just close the connection)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	
	// create a buffer to read data from the client
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return // client closed the connection or an error occurred
		}

		fmt.Printf("Received data: %s", string(buf[:n]))

		// to test the connection, we can send a simple response back to the client
		conn.Write([]byte("PONG\r\n"))
	}
}