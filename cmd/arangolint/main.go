// Package main runs the analyzer. It's the CLI entrypoint.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"go.augendre.info/arangolint/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer())
}
