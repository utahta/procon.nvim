package language

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/procon.nvim/src/cmp"
	"github.com/utahta/procon.nvim/src/pipeline"
)

type (
	Controller interface {
		SetFilename(filename string)
		Filename() string
		Sourcecode() (string, error)
		LanguageType() Type

		Valid() bool
		Build(ctx context.Context) error
		Compare(ctx context.Context, input string, output string) (string, error)
	}

	controller struct {
		file    *File
		builder Builder
	}
)

func NewController(rawLanguageType string) (Controller, error) {
	var b Builder
	switch strings.ToLower(strings.TrimSpace(rawLanguageType)) {
	case "c++14", "cpp14":
		b = newCpp14gccBuilder()
	case "c++", "c++17", "cpp17":
		b = newCpp17gccBuilder()
	default:
		return nil, errors.Errorf("%s is not supported for now", rawLanguageType)
	}

	var f *File
	switch b.Type() {
	case TypeCPP14gcc, TypeCPP17gcc:
		f = newFile(".cpp")
	default:
		f = newFile("")
	}

	return &controller{
		builder: b,
		file:    f,
	}, nil
}

func (c *controller) SetFilename(filename string) {
	c.file.SetName(filename)
}

func (c *controller) Filename() string {
	return c.file.Name()
}

func (c *controller) Sourcecode() (string, error) {
	return c.file.ReadAll()
}

func (c *controller) LanguageType() Type {
	return c.builder.Type()
}

func (c *controller) Valid() bool {
	if c.file.Name() == "" {
		return false
	}
	return true
}

func (c *controller) Build(ctx context.Context) error {
	return c.builder.Build(ctx, c.file.Name())
}

func (c *controller) Compare(ctx context.Context, input string, output string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	got, err := pipeline.Output(
		ctx,
		[]string{"echo", input},
		[]string{fmt.Sprintf("./%s", c.builder.ExecName())},
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to compare")
	}
	return cmp.Diff(output, strings.TrimSpace(string(got))), nil
}
