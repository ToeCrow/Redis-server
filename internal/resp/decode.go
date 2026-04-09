package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// ReadValue reads one RESP2 value from r.
func ReadValue(r io.Reader) (Value, error) {
	br := bufio.NewReader(r)
	return readValue(br)
}

func readValue(br *bufio.Reader) (Value, error) {
	b, err := br.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch b {
	case '+':
		s, err := readCRLFLine(br)
		if err != nil {
			return Value{}, err
		}
		return Value{Kind: KindSimpleString, Str: string(s)}, nil
	case '-':
		s, err := readCRLFLine(br)
		if err != nil {
			return Value{}, err
		}
		return Value{Kind: KindError, Str: string(s)}, nil
	case ':':
		s, err := readCRLFLine(br)
		if err != nil {
			return Value{}, err
		}
		n, err := strconv.ParseInt(string(s), 10, 64)
		if err != nil {
			return Value{}, fmt.Errorf("resp: invalid integer %q: %w", s, err)
		}
		return Value{Kind: KindInteger, Int: n}, nil
	case '$':
		return readBulk(br)
	case '*':
		return readArray(br)
	default:
		return Value{}, fmt.Errorf("resp: unknown type prefix %q", b)
	}
}

func readBulk(br *bufio.Reader) (Value, error) {
	line, err := readCRLFLine(br)
	if err != nil {
		return Value{}, err
	}
	n, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, fmt.Errorf("resp: invalid bulk length %q: %w", line, err)
	}
	if n == -1 {
		return Value{Kind: KindBulkString, BulkNull: true}, nil
	}
	if n < -1 {
		return Value{}, fmt.Errorf("resp: invalid bulk length %d", n)
	}
	body := make([]byte, n+2)
	if _, err := io.ReadFull(br, body); err != nil {
		return Value{}, err
	}
	if body[n] != '\r' || body[n+1] != '\n' {
		return Value{}, fmt.Errorf("resp: bulk string not followed by CRLF")
	}
	return Value{Kind: KindBulkString, Str: string(body[:n])}, nil
}

func readArray(br *bufio.Reader) (Value, error) {
	line, err := readCRLFLine(br)
	if err != nil {
		return Value{}, err
	}
	count, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, fmt.Errorf("resp: invalid array length %q: %w", line, err)
	}
	if count == -1 {
		return Value{Kind: KindArray, ArrayNull: true}, nil
	}
	if count < -1 {
		return Value{}, fmt.Errorf("resp: invalid array length %d", count)
	}
	elems := make([]Value, 0, count)
	for i := 0; i < count; i++ {
		v, err := readValue(br)
		if err != nil {
			return Value{}, err
		}
		elems = append(elems, v)
	}
	return Value{Kind: KindArray, Elems: elems}, nil
}

func readCRLFLine(br *bufio.Reader) ([]byte, error) {
	var line []byte
	for {
		b, err := br.ReadByte()
		if err != nil {
			if err == io.EOF && len(line) == 0 {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}
		if b == '\r' {
			next, err := br.ReadByte()
			if err != nil {
				return nil, err
			}
			if next != '\n' {
				return nil, fmt.Errorf("resp: expected LF after CR")
			}
			return line, nil
		}
		line = append(line, b)
	}
}
