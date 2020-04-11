package handler

import (
	"context"

	"github.com/utahta/procon.nvim/src/language"
)

type (
	Handler interface {
		Login(ctx context.Context) error
		Set(ctx context.Context) error
		Check(ctx context.Context, controller language.Controller) error
		Submit(ctx context.Context, controller language.Controller) error
		CurrentTask(ctx context.Context) string
	}
)
