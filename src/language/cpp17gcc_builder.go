package language

import (
	"context"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

type cpp17gccBuilder struct{}

func newCpp17gccBuilder() Builder {
	return &cpp17gccBuilder{}
}

func (b *cpp17gccBuilder) Build(ctx context.Context, filename string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "g++", "-std=gnu++17", "-Wall", "-Wextra", "-o", b.ExecName(), filename).Run(); err != nil {
		return errors.Wrapf(err, "failed to build: %s", filename)
	}
	return nil
}

func (b *cpp17gccBuilder) Type() Type {
	return TypeCPP17gcc
}

func (b *cpp17gccBuilder) ExecName() string {
	return "a.out"
}
