package atcoder

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/utahta/procon.nvim/src/procon-client/atcoder"
	"github.com/utahta/procon.nvim/src/vim/vimtest"
)

func TestHandler_Set(t *testing.T) {
	testcases := []struct {
		contest  string
		task     string
		wantTask string
	}{
		{
			contest:  "abc160",
			task:     "d",
			wantTask: "abc160_d",
		},
		{
			contest:  "dwango2015-prelims",
			task:     "1",
			wantTask: "dwango2015_prelims_1",
		},
		{
			contest:  "abc094",
			task:     "arc095_b",
			wantTask: "arc095_b",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.contest+"_"+tc.task, func(t *testing.T) {
			vim := vimtest.New()
			vim.CallStub = func(fname string, result interface{}, args ...interface{}) error {
				switch fname {
				case "input":
					if len(args) == 0 {
						return errors.New("invalid args")
					}
					prompt, ok := args[0].(string)
					if !ok {
						return errors.New("invalid type")
					}
					if strings.HasPrefix(prompt, "contest:") {
						*result.(*string) = tc.contest
					} else if strings.HasPrefix(prompt, "task:") {
						*result.(*string) = tc.task
					}
				}
				return nil
			}

			tempDir := os.TempDir()
			cli, err := atcoder.New(tempDir)
			if err != nil {
				t.Fatal(err)
			}
			h := NewHandler(vim, cli).(*Handler)
			if err := h.Set(context.Background()); err != nil {
				t.Fatal(err)
			}
			if want, got := tc.wantTask, h.task; want != got {
				t.Errorf("want:%v, got:%v", want, got)
			}
		})
	}
}
