package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"go.augendre.info/arangolint/pkg/analyzer"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc string
		dir  string
	}{
		{
			desc: "common",
			dir:  "common",
		},
		{
			desc: "cgo",
			dir:  "cgo",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc+"_"+test.dir, func(t *testing.T) {
			t.Parallel()

			anlzr := analyzer.NewAnalyzer()

			analysistest.Run(t, analysistest.TestData(), anlzr, test.dir)
		})
	}
}
