package language

import (
	"context"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

type cpp14gccBuilder struct{}

func newCpp14gccBuilder() Builder {
	return &cpp14gccBuilder{}
}

func (b *cpp14gccBuilder) Build(ctx context.Context, filename string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "g++", "-std=c++14", "-Wall", "-Wextra", "-o", b.ExecName(), filename).Run(); err != nil {
		return errors.Wrapf(err, "failed to build: %s", filename)
	}
	return nil
}

func (b *cpp14gccBuilder) Type() Type {
	return TypeCPP14gcc
}

func (b *cpp14gccBuilder) ExecName() string {
	return "a.out"
}
