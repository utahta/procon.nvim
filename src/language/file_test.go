package language

import (
	"testing"
)

func TestFile_SetName(t *testing.T) {
	testcases := []struct {
		name     string
		ext      string
		wantName string
	}{
		{
			name:     "abc160_a",
			ext:      ".cpp",
			wantName: "abc160_a.cpp",
		},
		{
			name:     "abc160_a.tmp",
			ext:      ".cpp",
			wantName: "abc160_a.tmp.cpp",
		},
		{
			name:     "abc160_a.cpp",
			ext:      ".cpp",
			wantName: "abc160_a.cpp",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFile(tc.ext)
			f.SetName(tc.name)

			if want, got := tc.wantName, f.Name(); want != got {
				t.Errorf("want %v, but got %v", want, got)
			}
		})
	}
}
