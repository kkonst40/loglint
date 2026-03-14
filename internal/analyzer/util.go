package analyzer

import (
	"go/ast"
	"go/types"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

func isLogCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if pass.TypesInfo == nil {
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

func reportWithSuggestedFix(pass *analysis.Pass, lit *ast.BasicLit, reportMsg, suggestMsg string, suggest string) {
	pass.Report(analysis.Diagnostic{
		Pos:     lit.Pos(),
		End:     lit.End(),
		Message: reportMsg,
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: suggestMsg,
				TextEdits: []analysis.TextEdit{
					{
						Pos:     lit.Pos(),
						End:     lit.End(),
						NewText: []byte(`"` + suggest + `"`),
					},
				},
			},
		},
	})
}

func startsWithUpper(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}
