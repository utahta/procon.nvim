package vim

import (
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
)

type (
	Vim interface {
		Call(fname string, result interface{}, args ...interface{}) error
		WriteOut(str string) error
		WriteErr(str string) error
		AC(no int) error
		WA(no int, diff string) error
	}

	proconVim struct {
		*nvim.Nvim
	}
)

func New(nvim *nvim.Nvim) Vim {
	return &proconVim{
		Nvim: nvim,
	}
}

func (c *proconVim) AC(no int) error {
	return c.Command(fmt.Sprintf(`:echohl ProconAC | echo "%d: AC\n" | echohl None`, no))
}

func (c *proconVim) WA(no int, diff string) error {
	outputs := []string{
		fmt.Sprintf(`:echohl ProconWA | echo "%d: WA\n" | echohl None`, no),
	}
	if diff != "" {
		outputs = append(outputs, fmt.Sprintf(`echo '%s'`, diff))
	}
	return c.Command(strings.Join(outputs, " | "))
}
