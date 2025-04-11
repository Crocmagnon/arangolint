package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Crocmagnon/arangolint/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer())
}
