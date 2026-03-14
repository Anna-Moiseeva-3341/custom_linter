package main

import (
	"github.com/Anna-Moiseeva-3341/custom_linter/pkg/loglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg := loglint.DefaultConfig()

	singlechecker.Main(loglint.NewAnalyzer(cfg))
}
