package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/procon.nvim/src/handler"
	"github.com/utahta/procon.nvim/src/handler/atcoder"
	"github.com/utahta/procon.nvim/src/language"
	atcoderclient "github.com/utahta/procon.nvim/src/procon-client/atcoder"
	"github.com/utahta/procon.nvim/src/vim"
)

type commander struct {
	vim       vim.Vim
	configDir string

	handler            handler.Handler
	languageController language.Controller
}

func newCommander(vim vim.Vim, configDir string) *commander {
	return &commander{
		vim:       vim,
		configDir: configDir,
	}
}

func (d *commander) Login() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h, err := d.ensureHandler()
	if err != nil {
		return err
	}
	return h.Login(ctx)
}

func (d *commander) Set() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h, err := d.ensureHandler()
	if err != nil {
		return err
	}
	if err := h.Set(ctx); err != nil {
		return err
	}

	c, err := d.ensureLanguageController()
	if err != nil {
		return err
	}

	c.SetFilename(h.CurrentTask(ctx))
	if _, err := os.Stat(c.Filename()); err == nil {
		// The file already exists, so enter any filename.
		filename := c.Filename()
		if err := d.vim.Call("input", &filename, "filename: ", filename); err != nil {
			return err
		}
		c.SetFilename(filename)
	}

	return d.vim.Call("execute", nil, fmt.Sprintf("edit %s", c.Filename()))
}

func (d *commander) Reset() error {
	d.handler = nil
	d.languageController = nil
	return d.Set()
}

func (d *commander) CurrentTask() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h, err := d.ensureHandler()
	if err != nil {
		return err
	}
	d.vim.WriteOut(fmt.Sprintf("%s\n", h.CurrentTask(ctx)))
	return nil
}

func (d *commander) Build() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := d.ensureLanguageController()
	if err != nil {
		return err
	}
	if !d.languageController.Valid() {
		if err := d.Set(); err != nil {
			return err
		}
	}

	if err := d.vim.Call("execute", nil, "write"); err != nil {
		return errors.Wrap(err, "failed to :w")
	}
	return c.Build(ctx)
}

func (d *commander) Check() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := d.Build(); err != nil {
		return err
	}
	h, err := d.ensureHandler()
	if err != nil {
		return err
	}
	return h.Check(ctx, d.languageController)
}

func (d *commander) Submit() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h, err := d.ensureHandler()
	if err != nil {
		return err
	}
	c, err := d.ensureLanguageController()
	if err != nil {
		return err
	}
	return h.Submit(ctx, c)
}

func (d *commander) ensureHandler() (handler.Handler, error) {
	if d.handler != nil {
		return d.handler, nil
	}
	var site string
	if err := d.vim.Call("input", &site, "site: "); err != nil {
		return nil, err
	}

	var h handler.Handler
	switch strings.ToLower(strings.TrimSpace(site)) {
	case "at", "atcoder":
		c, err := atcoderclient.New(d.configDir)
		if err != nil {
			return nil, err
		}
		h = atcoder.NewHandler(d.vim, c)
	default:
		return nil, errors.New("invalid site")
	}
	d.handler = h
	return d.handler, nil
}

func (d *commander) ensureLanguageController() (language.Controller, error) {
	if d.languageController != nil {
		return d.languageController, nil
	}
	var rawLanguageType string
	if err := d.vim.Call("input", &rawLanguageType, "language: "); err != nil {
		return nil, err
	}

	c, err := language.NewController(rawLanguageType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new language controller")
	}
	d.languageController = c
	return c, nil
}
