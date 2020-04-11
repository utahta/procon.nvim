package atcoder

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func runTestServer() *httptest.Server {
	mux := http.NewServeMux()

	copyTestHTML := func(w http.ResponseWriter, req *http.Request, filename string) {
		fp, err := os.Open(filename)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer fp.Close()
		io.Copy(w, fp)
	}

	mux.HandleFunc("/contests/abc160/tasks/abc160_d", func(w http.ResponseWriter, req *http.Request) {
		copyTestHTML(w, req, "testdata/abc160_d")
	})
	mux.HandleFunc("/contests/hitachi2020/tasks/hitachi2020_a", func(w http.ResponseWriter, req *http.Request) {
		copyTestHTML(w, req, "testdata/hitachi2020_a")
	})
	mux.HandleFunc("/contests/dwango2015-prelims/tasks/dwango2015_prelims_1", func(w http.ResponseWriter, req *http.Request) {
		copyTestHTML(w, req, "testdata/dwango2015_prelims_1")
	})
	mux.HandleFunc("/contests/arc011/tasks/arc011_1", func(w http.ResponseWriter, req *http.Request) {
		copyTestHTML(w, req, "testdata/arc011_1")
	})
	mux.HandleFunc("/contests/abc094/tasks/arc095_b", func(w http.ResponseWriter, req *http.Request) {
		copyTestHTML(w, req, "testdata/arc095_b")
	})
	mux.HandleFunc("/contests/arc046/tasks/arc046_a", func(w http.ResponseWriter, req *http.Request) {
		copyTestHTML(w, req, "testdata/arc046_a")
	})
	return httptest.NewServer(mux)
}

func TestClient_ListExamples(t *testing.T) {
	testcases := []struct {
		contest      string
		task         string
		wantExamples []Example
	}{
		{
			contest: "abc160",
			task:    "abc160_d",
			wantExamples: []Example{
				{Input: "5 2 4", Output: "5\n4\n1\n0"},
				{Input: "3 1 3", Output: "3\n0"},
				{Input: "7 3 7", Output: "7\n8\n4\n2\n0\n0"},
				{Input: "10 4 8", Output: "10\n12\n10\n8\n4\n1\n0\n0\n0"},
			},
		},
		{
			contest: "hitachi2020",
			task:    "hitachi2020_a",
			wantExamples: []Example{
				{Input: "hihi", Output: "Yes"},
				{Input: "hi", Output: "Yes"},
				{Input: "ha", Output: "No"},
			},
		},
		{
			contest: "dwango2015-prelims",
			task:    "dwango2015_prelims_1",
			wantExamples: []Example{
				{Input: "5\n2", Output: "2655"},
				{Input: "25\n25", Output: "13500"},
				{Input: "100\n1", Output: "52515"},
				{Input: "1\n1", Output: "540"},
			},
		},
		{
			contest: "arc011",
			task:    "arc011_1",
			wantExamples: []Example{
				{Input: "2 1 8", Output: "15"},
				{Input: "7 4 30", Output: "62"},
				{Input: "100 99 1000", Output: "90199"},
			},
		},
		{
			contest: "abc094",
			task:    "arc095_b",
			wantExamples: []Example{
				{Input: "5\n6 9 4 2 11", Output: "11 6"},
				{Input: "2\n100 0", Output: "100 0"},
			},
		},
		{
			contest: "arc046",
			task:    "arc046_a",
			wantExamples: []Example{
				{Input: "1", Output: "1"},
				{Input: "11", Output: "22"},
				{Input: "50", Output: "555555"},
			},
		},
	}

	ts := runTestServer()
	defer ts.Close()

	for _, tc := range testcases {
		t.Run(tc.contest+"_"+tc.task, func(t *testing.T) {
			c, err := New(os.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			c.URL, err = url.Parse(ts.URL)
			if err != nil {
				t.Fatal(err)
			}

			got, err := c.ListExamples(context.Background(), tc.contest, tc.task)
			if err != nil {
				t.Fatal(err)
			}

			if d := cmp.Diff(got, tc.wantExamples); d != "" {
				t.Fatal(d)
			}

			// for cache
			got, err = c.ListExamples(context.Background(), tc.contest, tc.task)
			if err != nil {
				t.Fatal(err)
			}
			if d := cmp.Diff(got, tc.wantExamples); d != "" {
				t.Fatal(d)
			}
		})
	}
}
