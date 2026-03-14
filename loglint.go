package loglint

import (
	"github.com/kkonst40/loglint/internal/analyzer"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

type LogLintPlugin struct {
	settings LogLintSettings
}

type LogLintSettings struct {
	CheckFirstChar       bool     `json:"check_first_char"`
	CheckNonEnglishChars bool     `json:"check_nonenglish_chars"`
	CheckSpecialChars    bool     `json:"check_special_chars"`
	CheckSensitiveWords  bool     `json:"check_sensitive_words"`
	SensitiveWords       []string `json:"sensitive_words"`
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[LogLintSettings](settings)
	if err != nil {
		return nil, err
	}

	return &LogLintPlugin{
		settings: s,
	}, nil
}

func (f *LogLintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	a := analyzer.New(
		f.settings.CheckFirstChar,
		f.settings.CheckNonEnglishChars,
		f.settings.CheckSpecialChars,
		f.settings.CheckSensitiveWords,
		f.settings.SensitiveWords,
	)

	return []*analysis.Analyzer{a}, nil
}

func (f *LogLintPlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
