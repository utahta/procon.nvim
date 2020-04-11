package language

import (
	"io/ioutil"
	"path/filepath"
)

type File struct {
	name string
	ext  string
}

func newFile(ext string) *File {
	return &File{
		ext: ext,
	}
}

func (f *File) SetName(name string) {
	var ext = filepath.Ext(name)
	if f.ext != ext {
		name += f.ext
	}
	f.name = name
}

func (f *File) Name() string {
	return f.name
}

func (f *File) ReadAll() (string, error) {
	bs, err := ioutil.ReadFile(f.name)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
