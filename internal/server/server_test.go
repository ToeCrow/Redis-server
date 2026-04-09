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

	go Serve(l, NewKVStore())

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
	go Serve(l, NewKVStore())

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
	go Serve(l, NewKVStore())

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
	go Serve(l, NewKVStore())

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

func TestSetReturnsOK(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	go Serve(l, NewKVStore())

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// SET foo bar
	req := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	if _, err := io.WriteString(conn, req); err != nil {
		t.Fatal(err)
	}
	br := bufio.NewReader(conn)
	line, err := br.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if line != "+OK\r\n" {
		t.Fatalf("got %q want \"+OK\\r\\n\"", line)
	}
}

func TestSetThenGetReturnsValue(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	go Serve(l, NewKVStore())

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	setReq := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	if _, err := io.WriteString(conn, setReq); err != nil {
		t.Fatal(err)
	}
	br := bufio.NewReader(conn)
	if line, err := br.ReadString('\n'); err != nil {
		t.Fatal(err)
	} else if line != "+OK\r\n" {
		t.Fatalf("SET: got %q want \"+OK\\r\\n\"", line)
	}

	getReq := "*2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n"
	if _, err := io.WriteString(conn, getReq); err != nil {
		t.Fatal(err)
	}
	out, err := br.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if out != "$3\r\n" {
		t.Fatalf("GET length line: got %q want \"$3\\r\\n\"", out)
	}
	body := make([]byte, 3+2)
	if _, err := io.ReadFull(br, body); err != nil {
		t.Fatal(err)
	}
	if string(body) != "bar\r\n" {
		t.Fatalf("GET bulk body: got %q", string(body))
	}
}

func TestGetMissingKeyReturnsNullBulk(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	go Serve(l, NewKVStore())

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	req := "*2\r\n$3\r\nGET\r\n$7\r\nmissing\r\n"
	if _, err := io.WriteString(conn, req); err != nil {
		t.Fatal(err)
	}
	br := bufio.NewReader(conn)
	line, err := br.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if line != "$-1\r\n" {
		t.Fatalf("got %q want \"$-1\\r\\n\"", line)
	}
}

func TestSetWrongArityReturnsError(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	go Serve(l, NewKVStore())

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// SET with one argument only (key missing value)
	if _, err := io.WriteString(conn, "*2\r\n$3\r\nSET\r\n$3\r\nfoo\r\n"); err != nil {
		t.Fatal(err)
	}
	br := bufio.NewReader(conn)
	line, err := br.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	want := "-ERR wrong number of arguments for 'set' command\r\n"
	if line != want {
		t.Fatalf("got %q want %q", line, want)
	}
}

func TestGetWrongArityReturnsError(t *testing.T) {
	t.Parallel()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	go Serve(l, NewKVStore())

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	if _, err := io.WriteString(conn, "*1\r\n$3\r\nGET\r\n"); err != nil {
		t.Fatal(err)
	}
	br := bufio.NewReader(conn)
	line, err := br.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	want := "-ERR wrong number of arguments for 'get' command\r\n"
	if line != want {
		t.Fatalf("got %q want %q", line, want)
	}
}
