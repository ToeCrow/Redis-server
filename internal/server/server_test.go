package server

import (
	"net"
	"testing"
	"time"
)

func TestServeAcceptsConnection(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()

	go Serve(l)

	time.Sleep(10 * time.Millisecond)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	conn.Close()
}
