package main

import (
	"os"
	"path/filepath"

	"github.com/neovim/go-client/nvim/plugin"
	"github.com/pkg/errors"
	"github.com/utahta/procon.nvim/src/vim"
)

func main() {
	plugin.Main(func(p *plugin.Plugin) error {
		configDir := os.Getenv("PROCON_NVIM_CONFIG_DIR")
		if configDir == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			configDir = filepath.Join(homeDir, ".config", "procon.nvim")
			if err := os.MkdirAll(filepath.Dir(configDir), os.ModePerm); err != nil {
				return errors.Wrap(err, "failed to make config directories")
			}
		}

		d := newCommander(vim.New(p.Nvim), configDir)
		p.HandleCommand(&plugin.CommandOptions{Name: "ProconLogin", NArgs: "0"}, d.Login)
		p.HandleCommand(&plugin.CommandOptions{Name: "ProconSet", NArgs: "0"}, d.Set)
		p.HandleCommand(&plugin.CommandOptions{Name: "ProconReset", NArgs: "0"}, d.Reset)
		p.HandleCommand(&plugin.CommandOptions{Name: "ProconCurrentTask", NArgs: "0"}, d.CurrentTask)
		p.HandleCommand(&plugin.CommandOptions{Name: "ProconBuild", NArgs: "0"}, d.Build)
		p.HandleCommand(&plugin.CommandOptions{Name: "ProconCheck", NArgs: "0"}, d.Check)
		p.HandleCommand(&plugin.CommandOptions{Name: "ProconSubmit", NArgs: "0"}, d.Submit)
		return nil
	})
}
