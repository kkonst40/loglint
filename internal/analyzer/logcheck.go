package analyzer

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

type logLinter struct {
	checkFirstChar       bool
	checkNonEnglishChars bool
	checkSpecialChars    bool
	checkSensitiveWords  bool
	sensitiveWords       []string
}

func New(checkFirstChar, checkNonEnglishChars, checkSpecialChars, checkSensitiveWords bool, sensitiveWords []string) *analysis.Analyzer {
	l := &logLinter{
		checkFirstChar:       checkFirstChar,
		checkNonEnglishChars: checkNonEnglishChars,
		checkSpecialChars:    checkSpecialChars,
		checkSensitiveWords:  checkSensitiveWords,
		sensitiveWords:       sensitiveWords,
	}

	analyzer := &analysis.Analyzer{
		Name:     "loglint",
		Doc:      "checks logging messages for style and security",
		Run:      l.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	return analyzer
}

func (l *logLinter) run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if isLogCall(pass, call) {
				l.checkLogMessage(pass, call)
			}
			return true
		})
	}
	return nil, nil
}

func (l *logLinter) checkLogMessage(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}

	firstArg := call.Args[0]

	if lit, ok := firstArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		msg := strings.Trim(lit.Value, "`\"")

		// starts with lowercase
		if l.checkFirstChar {
			if len(msg) > 0 && startsWithUpper(msg) {
				reportWithSuggestedFix(
					pass, lit,
					"log message should start with a lowercase letter",
					"make log message start with lowercase letter",
					l.correctMsg(msg),
				)
				return
			}
		}

		// english only
		if l.checkNonEnglishChars {
			for _, r := range msg {
				if unicode.IsLetter(r) && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
					reportWithSuggestedFix(
						pass, lit,
						"log message contains non-english characters",
						"make log message without non-english characters",
						l.correctMsg(msg),
					)

					return
				}
			}
		}

		// invalid characters
		if l.checkSpecialChars {
			for _, r := range msg {
				if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' {
					reportWithSuggestedFix(
						pass, lit,
						"log message contains invalid characters: emoji or special symbols",
						"make log message without invalid characters",
						l.correctMsg(msg),
					)

					return
				}
			}
		}
	}

	// sensitive data
	if l.checkSensitiveWords {
		l.checkSensitiveData(pass, call)
	}
}

func (l *logLinter) correctMsg(msg string) string {
	newMsgRunes := []rune{}

	for _, r := range msg {
		if l.checkNonEnglishChars {
			if unicode.IsLetter(r) && ((r > 'a' && r < 'z') || (r > 'A' && r < 'Z')) {
				newMsgRunes = append(newMsgRunes, r)
				continue
			}
		}

		if l.checkSpecialChars {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
				newMsgRunes = append(newMsgRunes, r)
			}
		}
	}

	if l.checkFirstChar && unicode.IsUpper(newMsgRunes[0]) {
		newMsgRunes[0] = unicode.ToLower(newMsgRunes[0])
	}

	return string(newMsgRunes)
}

func (l *logLinter) checkSensitiveData(pass *analysis.Pass, call *ast.CallExpr) {
	for _, arg := range call.Args {
		ast.Inspect(arg, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok {
				return true
			}

			lowerName := strings.ToLower(ident.Name)
			for _, word := range l.sensitiveWords {
				if strings.Contains(lowerName, word) {
					pass.Reportf(ident.Pos(), "potential sensitive data leak in log: variable '%s'", ident.Name)
					return false
				}
			}
			return true
		})
	}
}
