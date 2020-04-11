package language

import "context"

type Type int

type Builder interface {
	Build(ctx context.Context, filename string) error
	Type() Type
	ExecName() string
}

const (
	TypeCPP14gcc Type = iota
	TypeCPP17gcc
)
