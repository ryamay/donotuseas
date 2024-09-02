package donotuseas_test

import (
	"testing"

	"github.com/ryamay/donotuseas"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/findcall"
)

func init() {
	findcall.Analyzer.Flags.Set("arg", "myContext")
	findcall.Analyzer.Flags.Set("param", "context.Context")
}

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, donotuseas.Analyzer, "a")
}
