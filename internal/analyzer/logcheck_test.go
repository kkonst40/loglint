package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLogLint(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, New(true, true, true, true, []string{"password", "pass", "token", "apikey", "api_key"}), "a")
}
