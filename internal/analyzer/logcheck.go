package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglint",
	Doc:      "checks logging messages for style and security",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if isLogCall(pass, call) {
				checkLogMessage(pass, call)
			}
			return true
		})
	}
	return nil, nil
}

func checkLogMessage(pass *analysis.Pass, call *ast.CallExpr) {
	if len(call.Args) == 0 {
		return
	}

	firstArg := call.Args[0]

	if lit, ok := firstArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		msg := strings.Trim(lit.Value, "`\"")

		// starts with small
		if len(msg) > 0 && startsWithUpper(msg) {
			pass.Reportf(lit.Pos(), "log message should start with a lowercase letter")
			return
		}

		// english only
		for _, r := range msg {
			if unicode.IsLetter(r) && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
				pass.Reportf(lit.Pos(), "log message contains non-english characters")
				return
			}
		}

		// invalid characters
		for _, r := range msg {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' {
				pass.Reportf(lit.Pos(), "log message contains invalid characters: emoji or special symbols")
				return
			}
		}
	}

	// sensitive data
	checkSensitiveData(pass, call)
}

func startsWithUpper(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}

var sensitiveWords = []string{"password", "pass", "token", "secret", "apikey", "api_key"}

func checkSensitiveData(pass *analysis.Pass, call *ast.CallExpr) {
	for _, arg := range call.Args {
		ast.Inspect(arg, func(n ast.Node) bool {
			ident, ok := n.(*ast.Ident)
			if !ok {
				return true
			}

			lowerName := strings.ToLower(ident.Name)
			for _, word := range sensitiveWords {
				if strings.Contains(lowerName, word) {
					pass.Reportf(ident.Pos(), "potential sensitive data leak in log: variable '%s'", ident.Name)
					return false
				}
			}
			return true
		})
	}
}

func isLogCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.ObjectOf(sel.Sel)
	if obj == nil {
		return false
	}

	sig, ok := obj.Type().(*types.Signature)
	if !ok {
		return false
	}

	if sig.Recv() != nil {
		recvType := sig.Recv().Type().String()
		if recvType == "*log/slog.Logger" {
			return isLoggingMethod(sel.Sel.Name)
		}
		if recvType == "*go.uber.org/zap.Logger" || recvType == "*go.uber.org/zap.SugaredLogger" {
			return isLoggingMethod(sel.Sel.Name)
		}
	}

	if pkg := obj.Pkg(); pkg != nil {
		if pkg.Path() == "log/slog" || pkg.Path() == "go.uber.org/zap" {
			return isLoggingMethod(sel.Sel.Name)
		}
	}

	return false
}

func isLoggingMethod(name string) bool {
	switch name {
	case "Info", "Infof", "Infow",
		"Error", "Errorf", "Errorw",
		"Warn", "Warnf", "Warnw",
		"Debug", "Debugf", "Debugw",
		"Fatal", "Fatalf", "Fatalw",
		"Panic", "Panicf", "Panicw":
		return true
	}
	return false
}
