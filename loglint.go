package loglint

import (
	"github.com/kkonst40/loglint/internal/analyzer"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

type LogLintPlugin struct{}

func New(settings any) (register.LinterPlugin, error) {
	return &LogLintPlugin{}, nil
}

func (f *LogLintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}, nil
}

func (f *LogLintPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
