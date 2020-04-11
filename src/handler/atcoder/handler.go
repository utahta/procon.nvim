package atcoder

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/utahta/procon.nvim/src/handler"
	"github.com/utahta/procon.nvim/src/language"
	"github.com/utahta/procon.nvim/src/procon-client/atcoder"
	"github.com/utahta/procon.nvim/src/vim"
)

type Handler struct {
	vim    vim.Vim
	client *atcoder.Client

	contest string
	task    string
}

func NewHandler(vim vim.Vim, client *atcoder.Client) handler.Handler {
	return &Handler{
		vim:    vim,
		client: client,
	}
}

func (h *Handler) Login(ctx context.Context) error {
	username := os.Getenv("PROCON_NVIM_ATCODER_USERNAME")
	password := os.Getenv("PROCON_NVIM_ATCODER_PASSWORD")
	if username == "" {
		if err := h.vim.Call("input", &username, "username: "); err != nil {
			return err
		}
	}
	if password == "" {
		if err := h.vim.Call("inputsecret", &password, "password: "); err != nil {
			return err
		}
	}

	if err := h.client.Login(ctx, username, password); err != nil {
		return err
	}
	return nil
}

func (h *Handler) Set(_ context.Context) error {
	if err := h.vim.Call("input", &h.contest, "contest: ", h.contest); err != nil {
		return err
	}
	if err := h.vim.Call("input", &h.task, "task: "); err != nil {
		return err
	}
	h.canonicalizeTask()
	return nil
}

func (h *Handler) Check(ctx context.Context, c language.Controller) error {
	examples, err := h.client.ListExamples(ctx, h.contest, h.task)
	if err != nil {
		return errors.Wrap(err, "failed to get examples")
	}

	type result struct {
		no   int
		diff string
	}
	var results []result
	for i, e := range examples {
		diff, err := c.Compare(ctx, e.Input, e.Output)
		if err != nil {
			return errors.Wrap(err, "failed to check")
		}
		results = append(results, result{
			no:   i + 1,
			diff: diff,
		})
	}

	for _, res := range results {
		if res.diff == "" {
			h.vim.AC(res.no)
		} else {
			h.vim.WA(res.no, res.diff)
		}
	}
	return nil
}

func (h *Handler) Submit(ctx context.Context, c language.Controller) error {
	if !h.client.IsLoggedIn() {
		if err := h.Login(ctx); err != nil {
			return err
		}
	}

	var languageId string
	switch c.LanguageType() {
	case language.TypeCPP14gcc:
		languageId = atcoder.LanguageIDCPP14gcc
	case language.TypeCPP17gcc:
		languageId = atcoder.LanguageIDCPP17gcc
	default:
		return errors.Errorf("%d language is not supported yet in this plugin.", c.LanguageType())
	}

	sourcecode, err := c.Sourcecode()
	if err != nil {
		return errors.Wrap(err, "failed to read source code")
	}

	redirectURL, err := h.client.Submit(ctx, h.contest, h.task, languageId, sourcecode)
	if err != nil {
		return errors.Wrap(err, "failed to submit")
	}

	if err := browser.OpenURL(redirectURL); err != nil {
		return errors.Wrap(err, "failed to open browser")
	}
	return nil
}

func (h *Handler) CurrentTask(_ context.Context) string {
	return h.task
}

func (h *Handler) canonicalizeTask() {
	if strings.Contains(h.task, "_") {
		return
	}
	h.task = fmt.Sprintf("%s_%s", strings.ReplaceAll(h.contest, "-", "_"), h.task)
}
