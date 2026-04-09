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
	Serve(l, NewKVStore())
	return nil
}

// Serve accepts connections on l until l is closed. kv is shared across clients.
func Serve(l net.Listener, kv *KVStore) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, kv)
	}
}

func handleConnection(conn net.Conn, kv *KVStore) {
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
		out := dispatchRequest(req, kv)
		if err := resp.WriteValue(conn, out); err != nil {
			return
		}
	}
}

func dispatchRequest(v resp.Value, kv *KVStore) resp.Value {
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
	case "SET":
		if len(v.Elems) != 3 {
			return resp.Err("ERR wrong number of arguments for 'set' command")
		}
		keyArg := v.Elems[1]
		valArg := v.Elems[2]
		if keyArg.Kind != resp.KindBulkString || keyArg.BulkNull {
			return resp.Err("ERR value is not a valid bulk string")
		}
		if valArg.Kind != resp.KindBulkString || valArg.BulkNull {
			return resp.Err("ERR value is not a valid bulk string")
		}
		kv.Set(keyArg.Str, valArg.Str)
		return resp.Simple("OK")
	case "GET":
		if len(v.Elems) != 2 {
			return resp.Err("ERR wrong number of arguments for 'get' command")
		}
		keyArg := v.Elems[1]
		if keyArg.Kind != resp.KindBulkString || keyArg.BulkNull {
			return resp.Err("ERR value is not a valid bulk string")
		}
		val, ok := kv.Get(keyArg.Str)
		if !ok {
			return resp.Bulk("", true)
		}
		return resp.Bulk(val, false)
	default:
		return resp.Err("ERR unknown command '" + cmd.Str + "'")
	}
}
