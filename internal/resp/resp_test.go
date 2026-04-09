package resp

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestDecodeReadmeExamples(t *testing.T) {
	t.Parallel()

	t.Run("ping array", func(t *testing.T) {
		t.Parallel()
		raw := "*1\r\n$4\r\nping\r\n"
		v, err := ReadValue(strings.NewReader(raw))
		if err != nil {
			t.Fatal(err)
		}
		want := Value{
			Kind: KindArray,
			Elems: []Value{
				{Kind: KindBulkString, Str: "ping"},
			},
		}
		if !equalValue(v, want) {
			t.Fatalf("got %+v want %+v", v, want)
		}
	})

	t.Run("simple OK", func(t *testing.T) {
		t.Parallel()
		v, err := ReadValue(strings.NewReader("+OK\r\n"))
		if err != nil {
			t.Fatal(err)
		}
		want := Simple("OK")
		if !equalValue(v, want) {
			t.Fatalf("got %+v want %+v", v, want)
		}
	})

	t.Run("null bulk", func(t *testing.T) {
		t.Parallel()
		v, err := ReadValue(strings.NewReader("$-1\r\n"))
		if err != nil {
			t.Fatal(err)
		}
		want := Bulk("", true)
		if !equalValue(v, want) {
			t.Fatalf("got %+v want %+v", v, want)
		}
	})
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		v    Value
	}{
		{"simple", Simple("OK")},
		{"error", Err("ERR unknown")},
		{"int zero", Integer(0)},
		{"int neg", Integer(-99)},
		{"bulk empty", Bulk("", false)},
		{"bulk null", Bulk("", true)},
		{"bulk data", Bulk("hello", false)},
		{"array empty", Array([]Value{}, false)},
		{"array null", Array(nil, true)},
		{"array one bulk", Array([]Value{{Kind: KindBulkString, Str: "ping"}}, false)},
		{"nested", Array([]Value{
			{Kind: KindArray, Elems: []Value{Integer(1), Simple("x")}},
		}, false)},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			if err := WriteValue(&buf, tc.v); err != nil {
				t.Fatal(err)
			}
			got, err := ReadValue(&buf)
			if err != nil {
				t.Fatal(err)
			}
			if !equalValue(got, tc.v) {
				t.Fatalf("round trip\ngot %+v\nwant %+v", got, tc.v)
			}
		})
	}
}

func TestEmptyBulkString(t *testing.T) {
	t.Parallel()
	raw := "$0\r\n\r\n"
	v, err := ReadValue(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	want := Bulk("", false)
	if !equalValue(v, want) {
		t.Fatalf("got %+v want %+v", v, want)
	}
}

func TestReadValueFromSequential(t *testing.T) {
	t.Parallel()
	// Two back-to-back RESP values; a single bufio.Reader must preserve both.
	raw := "+OK\r\n$5\r\nhello\r\n"
	br := bufio.NewReader(strings.NewReader(raw))
	v1, err := ReadValueFrom(br)
	if err != nil {
		t.Fatal(err)
	}
	if !equalValue(v1, Simple("OK")) {
		t.Fatalf("first: got %+v", v1)
	}
	v2, err := ReadValueFrom(br)
	if err != nil {
		t.Fatal(err)
	}
	if !equalValue(v2, Bulk("hello", false)) {
		t.Fatalf("second: got %+v", v2)
	}
}

func TestNestedArray(t *testing.T) {
	t.Parallel()
	// *2\r\n*1\r\n:42\r\n$1\r\nx\r\n
	raw := "*2\r\n*1\r\n:42\r\n$1\r\nx\r\n"
	v, err := ReadValue(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	want := Array([]Value{
		Array([]Value{Integer(42)}, false),
		Bulk("x", false),
	}, false)
	if !equalValue(v, want) {
		t.Fatalf("got %+v want %+v", v, want)
	}
}

func equalValue(a, b Value) bool {
	if a.Kind != b.Kind {
		return false
	}
	switch a.Kind {
	case KindSimpleString, KindError:
		return a.Str == b.Str
	case KindInteger:
		return a.Int == b.Int
	case KindBulkString:
		if a.BulkNull != b.BulkNull {
			return false
		}
		if a.BulkNull {
			return true
		}
		return a.Str == b.Str
	case KindArray:
		if a.ArrayNull != b.ArrayNull {
			return false
		}
		if a.ArrayNull {
			return true
		}
		if len(a.Elems) != len(b.Elems) {
			return false
		}
		for i := range a.Elems {
			if !equalValue(a.Elems[i], b.Elems[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
