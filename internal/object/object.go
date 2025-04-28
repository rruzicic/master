package object

import (
	"bytes"
	"fmt"
	"interpreter/internal/ast"
	"strings"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	ARRAY_OBJ        = "ARRAY"
	STDFUNC_OBJ      = "STDFUNC"
)

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Float struct {
	Value float64
}

func (i *Float) Inspect() string  { return fmt.Sprintf("%f", i.Value) }
func (i *Float) Type() ObjectType { return FLOAT_OBJ }

type Boolean struct {
	Value bool
}

func (i *Boolean) Inspect() string  { return fmt.Sprintf("%t", i.Value) }
func (i *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type Null struct{}

func (i *Null) Inspect() string  { return "null" }
func (i *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return fmt.Sprintf("%v", rv.Value) }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Error struct {
	Error string
}

func (e *Error) Inspect() string  { return fmt.Sprintf("ERROR: %s", e.Error) }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

type Function struct {
	Params []ast.IdentifierExpression
	Body   *ast.BlockStatement
	Env    *Environment
}

func (e *Function) Inspect() string  { return "<fn>" }
func (e *Function) Type() ObjectType { return FUNCTION_OBJ }

type String struct {
	Value string
}

func (i *String) Inspect() string  { return i.Value }
func (i *String) Type() ObjectType { return STRING_OBJ }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type StdFunction struct {
	Fun func(args ...Object) Object
}

func (sf *StdFunction) Type() ObjectType { return STDFUNC_OBJ }
func (sf *StdFunction) Inspect() string  { return "<std fun>" }
