package resp

import (
	"fmt"
	"io"
)

// WriteValue writes v in RESP2 form to w.
func WriteValue(w io.Writer, v Value) error {
	switch v.Kind {
	case KindSimpleString:
		if _, err := fmt.Fprintf(w, "+%s\r\n", v.Str); err != nil {
			return err
		}
		return nil
	case KindError:
		if _, err := fmt.Fprintf(w, "-%s\r\n", v.Str); err != nil {
			return err
		}
		return nil
	case KindInteger:
		if _, err := fmt.Fprintf(w, ":%d\r\n", v.Int); err != nil {
			return err
		}
		return nil
	case KindBulkString:
		if v.BulkNull {
			_, err := io.WriteString(w, "$-1\r\n")
			return err
		}
		b := v.Str
		if _, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(b), b); err != nil {
			return err
		}
		return nil
	case KindArray:
		if v.ArrayNull {
			_, err := io.WriteString(w, "*-1\r\n")
			return err
		}
		elems := v.Elems
		if _, err := fmt.Fprintf(w, "*%d\r\n", len(elems)); err != nil {
			return err
		}
		for i := range elems {
			if err := WriteValue(w, elems[i]); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("resp: unknown value kind %d", v.Kind)
	}
}
