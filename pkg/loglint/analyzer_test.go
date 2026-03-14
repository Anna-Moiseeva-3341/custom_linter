package loglint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLogLint(t *testing.T) {
	testdata := analysistest.TestData()
	cfg := DefaultConfig()
	analysistest.Run(t, testdata, NewAnalyzer(cfg), "logtests")
}
