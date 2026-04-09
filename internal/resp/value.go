package resp

// Kind identifies which RESP2 value is held in Value.
type Kind int

const (
	KindSimpleString Kind = iota
	KindError
	KindInteger
	KindBulkString
	KindArray
)

// Value is a RESP2-encoded value (simple string, error, integer, bulk string, or array).
// Only the fields relevant to Kind are meaningful:
//   - KindSimpleString, KindError: Str
//   - KindInteger: Int
//   - KindBulkString: BulkNull and Str (Str ignored when BulkNull is true)
//   - KindArray: ArrayNull and Elems (null array vs empty vs elements)
type Value struct {
	Kind Kind

	Str       string
	Int       int64
	BulkNull  bool
	Elems     []Value
	ArrayNull bool
}

// Simple returns a simple string value (+...).
func Simple(s string) Value {
	return Value{Kind: KindSimpleString, Str: s}
}

// Err returns an error value (-...).
func Err(msg string) Value {
	return Value{Kind: KindError, Str: msg}
}

// Integer returns an integer value (:...).
func Integer(n int64) Value {
	return Value{Kind: KindInteger, Int: n}
}

// Bulk returns a bulk string; pass bulkNull true for Redis null bulk ($-1).
func Bulk(s string, bulkNull bool) Value {
	return Value{Kind: KindBulkString, Str: s, BulkNull: bulkNull}
}

// Array returns an array value. Use ArrayNull for *-1; empty slice for *0.
func Array(elems []Value, arrayNull bool) Value {
	return Value{Kind: KindArray, Elems: elems, ArrayNull: arrayNull}
}
