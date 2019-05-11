package object

import "fmt"

// Type is a constant representing the type of an underlying object.
type Type string

// These are all the possible value types.
const (
	IntegerType = "INTEGER"
	BooleanType = "BOOLEAN"
	NullType    = "NULL"
	FuncType    = "FUNCTION"
	ErrType     = "ERROR"
)

// Object is an interface of any object in WELP.
type Object interface {
	Type() Type
	Inspect() string
}

// Integer represents WELP's integer values.
type Integer struct {
	Value int64
}

// Type implements Object.
func (i *Integer) Type() Type {
	return IntegerType
}

// Inspect implements Object.
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Boolean represents WELP's bool values.
type Boolean struct {
	Value bool
}

// Type implements Object.
func (b *Boolean) Type() Type {
	return BooleanType
}

// Inspect implements Object.
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}

// Null represents WELP's null values.
type Null struct {
}

// Type implements Object.
func (n *Null) Type() Type {
	return NullType
}

// Inspect implements Object.
func (n *Null) Inspect() string {
	return "null"
}

// Func represents WELP's function values.
type Func struct {
	Name string
}

// Type implements Object.
func (f *Func) Type() Type {
	return FuncType
}

// Inspect implements Object.
func (f *Func) Inspect() string {
	return fmt.Sprintf("<func %s>", f.Name)
}

// Error represents WELP's error values.
type Error struct {
	Err error
}

// Type implements Object.
func (e *Error) Type() Type {
	return ErrType
}

// Inspect implements Object.
func (e *Error) Inspect() string {
	return fmt.Sprintf("ERR: %v", e.Err)
}
