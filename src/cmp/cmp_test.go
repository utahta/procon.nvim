package cmp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Diff(t *testing.T) {
	testcases := []struct {
		name string
		x    string
		y    string
		want string
	}{
		{
			name: "no diff",
			x:    "aaa",
			y:    "aaa",
			want: "",
		},
		{
			name: "simple",
			x:    "aaa",
			y:    "bbb",
			want: "-\naaa\n+\nbbb\n",
		},
		{
			name: "newline",
			x:    "aaa\naaa\naaa",
			y:    "aaa\naaa\nbbb",
			want: "-\naaa\naaa\naaa\n+\naaa\naaa\nbbb\n",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := Diff(tc.x, tc.y)
			if d := cmp.Diff(tc.want, got); d != "" {
				t.Fatal(d)
			}
		})
	}
}
