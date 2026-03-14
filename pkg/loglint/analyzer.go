package loglint

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/types/typeutil"
)

type Config struct {
	EnabledChecks     map[string]bool
	ForbiddenWords    []string
	ForbiddenSymbols  []string
	ForbiddenPatterns []string
}

type LogLinter struct {
	cfg              Config
	forbiddenSymbols string
	sensitiveRegex   *regexp.Regexp
}

func DefaultConfig() Config {
	return Config{
		EnabledChecks:     map[string]bool{"lowercase": true, "language": true, "emoji": true, "symbols": true, "sensitive": true},
		ForbiddenSymbols:  []string{"!", "?", ";"},
		ForbiddenPatterns: []string{"..."},
		ForbiddenWords:    []string{"password", "token", "pass"},
	}
}

func NewAnalyzer(cfg Config) *analysis.Analyzer {
	var symbols []string
	var reg *regexp.Regexp

	for _, s := range cfg.ForbiddenSymbols {
		if len(s) == 1 {
			symbols = append(symbols, s)
		}
	}

	if len(cfg.ForbiddenWords) == 0 {
		reg = regexp.MustCompile("(?i)^$")
	} else {
		reg = regexp.MustCompile(`(?i)(` + strings.Join(cfg.ForbiddenWords, "|") + `)\s*[:=]`)
	}

	l := &LogLinter{
		cfg:              cfg,
		forbiddenSymbols: strings.Join(symbols, ""),
		sensitiveRegex:   reg,
	}

	return &analysis.Analyzer{
		Name:     "customloglint",
		Doc:      "checks log messages for formatting, language, sensitive data",
		Run:      l.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func (l *LogLinter) run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) == 0 {
				return true
			}

			fn, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
			if !ok || fn.Pkg() == nil {
				return true
			}

			if l.cfg.EnabledChecks["sensitive"] {
				l.checkSensitiveData(pass, call)
			}

			msg := getMsg(pass, call)
			if msg == nil {
				return true
			}

			tv := pass.TypesInfo.Types[msg]
			if tv.Value == nil || tv.Value.Kind() != constant.String {
				return true
			}

			msgStr := constant.StringVal(tv.Value)
			if len(msgStr) == 0 {
				return true
			}

			if l.cfg.EnabledChecks["lowercase"] && !checkLowerCase(msgStr) {
				pass.Reportf(msg.Pos(), "log messages must start with lowercase letter")
			}

			if l.cfg.EnabledChecks["language"] && !checkLanguage(msgStr) {
				pass.Reportf(msg.Pos(), "log messages must contain only english")
			}

			if l.cfg.EnabledChecks["emoji"] && !checkEmoji(msgStr) {
				pass.Reportf(msg.Pos(), "log messages must not contain emoji")
			}

			if l.cfg.EnabledChecks["symbols"] && !l.checkSymbols(msgStr) {
				pass.Reportf(msg.Pos(), "log messages must not contain special symbols")
			}
			return true
		})
	}
	return nil, nil
}

func isLogPkg(pkg *types.Package) bool {
	switch pkg.Path() {
	case "log/slog", "go.uber.org/zap", "log":
		return true
	default:
		return false
	}
}

func isLogMethod(name string) bool {
	switch name {
	case "Info", "Infof", "Infow", "Warn", "Warnf", "Warnw",
		"Error", "Errorf", "Errorw", "Debug", "Debugf",
		"Debugw", "Fatal", "Fatalf", "Fatalw", "Print", "Printf", "Println":
		return true
	default:
		return false
	}
}

func msgParamName(name string) bool {
	switch name {
	case "msg", "message", "format", "s", "template", "args", "v":
		return true
	default:
		return false
	}
}

func getMsg(p *analysis.Pass, c *ast.CallExpr) ast.Expr {
	fn, ok := typeutil.Callee(p.TypesInfo, c).(*types.Func)
	if !ok || !isLogPkg(fn.Pkg()) {
		return nil
	}

	if !isLogMethod(fn.Name()) {
		return nil
	}

	signature, ok := fn.Type().(*types.Signature)
	if !ok {
		return nil
	}

	params := signature.Params()
	msgId := -1
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		if msgParamName(p.Name()) {
			msgId = i
			break
		}
	}

	if msgId != -1 && len(c.Args) > msgId {
		return c.Args[msgId]
	}

	for _, arg := range c.Args {
		tv := p.TypesInfo.Types[arg]
		if tv.IsValue() && tv.Value != nil && tv.Value.Kind() == constant.String {
			return arg
		}
	}
	return nil
}

func checkLowerCase(msg string) bool {
	firstRune := []rune(msg)[0]
	return !unicode.IsUpper(firstRune)
}

func checkLanguage(msg string) bool {
	for _, s := range msg {
		if s > unicode.MaxASCII && unicode.IsLetter(s) {
			return false
		}
	}
	return true
}

func checkEmoji(msg string) bool {
	for _, s := range msg {
		if s > unicode.MaxASCII && !unicode.IsLetter(s) {
			return false
		}
	}
	return true
}

func (l *LogLinter) checkSymbols(msg string) bool {
	if l.forbiddenSymbols == "" {
		return true
	}

	if strings.ContainsAny(msg, l.forbiddenSymbols) {
		return false
	}

	for _, pattern := range l.cfg.ForbiddenPatterns {
		if strings.Contains(msg, pattern) {
			return false
		}
	}
	return true
}

func (l *LogLinter) isSensitiveWord(msg string) bool {
	msg = strings.ToLower(msg)
	for _, word := range l.cfg.ForbiddenWords {
		if strings.Contains(msg, strings.ToLower(word)) {
			return true
		}
	}

	return false
}

func (l *LogLinter) checkSensitiveData(p *analysis.Pass, call *ast.CallExpr) {
	for _, arg := range call.Args {
		if l.checkExpr(p, arg) {
			return
		}
	}
}

func (l *LogLinter) checkExpr(p *analysis.Pass, expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		if l.isSensitiveWord(e.Name) {
			p.Reportf(e.Pos(), "log messages must not contain sensitive data")
			return true
		}
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			if l.sensitiveRegex.MatchString(e.Value) {
				p.Reportf(e.Pos(), "log messages must not contain sensitive data")
				return true
			}
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			if l.checkExpr(p, e.X) {
				return true
			}
			return l.checkExpr(p, e.Y)
		}
	case *ast.CallExpr:
		for _, arg := range e.Args {
			if l.checkExpr(p, arg) {
				return true
			}
		}
	}
	return false
}
