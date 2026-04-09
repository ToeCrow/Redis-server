package server

import (
	"fmt"
	"net"
)

// Run listens on addr and serves clients until the listener is closed.
func Run(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	fmt.Println("Redis-clone listen to", addr)
	Serve(l)
	return nil
}

// Serve accepts connections on l until l is closed.
func Serve(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		fmt.Printf("Received data: %s", string(buf[:n]))

		conn.Write([]byte("PONG\r\n"))
	}
}
