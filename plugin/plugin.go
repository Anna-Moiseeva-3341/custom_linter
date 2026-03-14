package plugin

import (
	"fmt"

	"github.com/Anna-Moiseeva-3341/custom_linter/pkg/loglint"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("customloglint", New)
}

func New(conf any) (register.LinterPlugin, error) {
	settings, ok := conf.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid config type")
	}

	cfg := loglint.Config{
		ForbiddenWords:   parseStrings(settings["forbidden_words"]),
		ForbiddenSymbols: parseStrings(settings["forbidden_symbols"]),
		EnabledChecks:    parseBoolMap(settings["enabled_checks"]),
	}
	return &customPlugin{cfg: cfg}, nil
}

type customPlugin struct {
	cfg loglint.Config
}

func (p *customPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		loglint.NewAnalyzer(p.cfg),
	}, nil
}

func parseStrings(in interface{}) []string {
	var res []string
	if list, ok := in.([]interface{}); ok {
		for _, v := range list {
			res = append(res, v.(string))
		}
	}
	return res
}

func parseBoolMap(in interface{}) map[string]bool {
	res := make(map[string]bool)
	if m, ok := in.(map[string]interface{}); ok {
		for k, v := range m {
			res[k] = v.(bool)
		}
	}
	return res
}

func (p *customPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
