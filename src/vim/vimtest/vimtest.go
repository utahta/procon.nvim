package vimtest

import "github.com/utahta/procon.nvim/src/vim"

type Vim struct {
	CallStub func(fname string, result interface{}, args ...interface{}) error
}

var _ vim.Vim = (*Vim)(nil)

func New() *Vim {
	return &Vim{
		CallStub: func(_ string, _ interface{}, _ ...interface{}) error { return nil },
	}
}

func (v *Vim) Call(fname string, result interface{}, args ...interface{}) error {
	return v.CallStub(fname, result, args...)
}

func (v *Vim) WriteOut(str string) error {
	return nil
}

func (v *Vim) WriteErr(str string) error {
	return nil
}

func (v *Vim) AC(no int) error {
	return nil
}

func (v *Vim) WA(no int, diff string) error {
	return nil
}
