package cmp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type diffReporter struct {
	path  cmp.Path
	diffs []string
}

func (r *diffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *diffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		// vx and vy should be string
		vx, vy := r.path.Last().Values()
		if vx.Kind() == reflect.String && vy.Kind() == reflect.String {
			r.diffs = append(r.diffs, fmt.Sprintf("-\n%s\n+\n%s\n", vx, vy))
		} else {
			// maybe does not reach here.
			r.diffs = append(r.diffs, fmt.Sprintf("%v\n\t- %v\n\t+ %v\n", r.path, vx, vy))
		}
	}
}

func (r *diffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *diffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}

func Diff(x, y string) string {
	r := &diffReporter{}
	if cmp.Equal(x, y, cmp.Reporter(r)) {
		return ""
	}
	return r.String()
}
