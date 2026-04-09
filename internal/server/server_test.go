package server

import (
	"bufio"
	"io"
	"net"
	"testing"
	"time"
)

func TestServeAcceptsConnection(t *testing.T) {
	t.Parallel()
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

func TestPingReturnsPong(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	go Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// *1\r\n$4\r\nping\r\n
	if _, err := io.WriteString(conn, "*1\r\n$4\r\nping\r\n"); err != nil {
		t.Fatal(err)
	}
	br := bufio.NewReader(conn)
	line, err := br.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if line != "+PONG\r\n" {
		t.Fatalf("got %q want \"+PONG\\r\\n\"", line)
	}
}

func TestEchoReturnsBulkString(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	go Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// *2\r\n$4\r\nECHO\r\n$11\r\nHello World\r\n
	req := "*2\r\n$4\r\nECHO\r\n$11\r\nHello World\r\n"
	if _, err := io.WriteString(conn, req); err != nil {
		t.Fatal(err)
	}
	br := bufio.NewReader(conn)
	out, err := br.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if out != "$11\r\n" {
		t.Fatalf("length line: got %q want \"$11\\r\\n\"", out)
	}
	body := make([]byte, 11+2)
	if _, err := io.ReadFull(br, body); err != nil {
		t.Fatal(err)
	}
	if string(body) != "Hello World\r\n" {
		t.Fatalf("bulk body: got %q", string(body))
	}
}

func TestUnknownCommandReturnsError(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	go Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	if _, err := io.WriteString(conn, "*1\r\n$7\r\nUNKNOWN\r\n"); err != nil {
		t.Fatal(err)
	}
	br := bufio.NewReader(conn)
	line, err := br.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if line != "-ERR unknown command 'UNKNOWN'\r\n" {
		t.Fatalf("got %q", line)
	}
}
