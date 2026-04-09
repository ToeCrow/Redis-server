package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/thokro/redis-server/internal/resp"
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

	br := bufio.NewReader(conn)
	for {
		req, err := resp.ReadValueFrom(br)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			_ = resp.WriteValue(conn, resp.Err("ERR Protocol error: "+err.Error()))
			return
		}
		out := dispatchRequest(req)
		if err := resp.WriteValue(conn, out); err != nil {
			return
		}
	}
}

func dispatchRequest(v resp.Value) resp.Value {
	if v.Kind != resp.KindArray || v.ArrayNull {
		return resp.Err("ERR Protocol error: expected array")
	}
	if len(v.Elems) == 0 {
		return resp.Err("ERR Protocol error: empty array")
	}
	cmd := v.Elems[0]
	if cmd.Kind != resp.KindBulkString || cmd.BulkNull {
		return resp.Err("ERR Protocol error: command must be a bulk string")
	}
	name := strings.ToUpper(cmd.Str)
	switch name {
	case "PING":
		switch len(v.Elems) {
		case 1:
			return resp.Simple("PONG")
		case 2:
			arg := v.Elems[1]
			if arg.Kind != resp.KindBulkString || arg.BulkNull {
				return resp.Err("ERR value is not a valid bulk string")
			}
			return resp.Bulk(arg.Str, false)
		default:
			return resp.Err("ERR wrong number of arguments for 'ping' command")
		}
	case "ECHO":
		if len(v.Elems) != 2 {
			return resp.Err("ERR wrong number of arguments for 'echo' command")
		}
		arg := v.Elems[1]
		if arg.Kind != resp.KindBulkString || arg.BulkNull {
			return resp.Err("ERR value is not a valid bulk string")
		}
		return resp.Bulk(arg.Str, false)
	default:
		return resp.Err("ERR unknown command '" + cmd.Str + "'")
	}
}
