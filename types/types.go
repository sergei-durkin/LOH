package types

import (
	"strings"
)

type BasicType int64

const (
	Uint BasicType = iota
	Uint8
	Uint16
	Uint32
	Uint64

	Int
	Int8
	Int16
	Int32
	Int64

	Pointer
	Bool
)

type TypeInfo interface {
	Size() int64
}

func (t BasicType) Size() int64 {
	return size[t]
}

var size = map[BasicType]int64{
	Uint:   4,
	Uint8:  1,
	Uint16: 2,
	Uint32: 4,
	Uint64: 8,

	Int:   4,
	Int8:  1,
	Int16: 2,
	Int32: 4,
	Int64: 8,

	Pointer: 8,
	Bool:    1,
}

var types = map[string]TypeInfo{
	"byte": Uint8,
	"bool": Uint8,

	"uint":   Uint,
	"uint8":  Uint8,
	"uint16": Uint16,
	"uint32": Uint32,
	"uint64": Uint64,

	"int":   Int,
	"int8":  Int8,
	"int16": Int16,
	"int32": Int32,
	"int64": Int64,

	"pointer": Pointer,
	"ptr":     Pointer,
}

func Info(t string) TypeInfo {
	return types[strings.ToLower(t)]
}
