// Package plugin is intended to be imported by a custom
// golangci-lint config as described here: https://golangci-lint.run/plugins/module-plugins/
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"go.augendre.info/arangolint/pkg/analyzer"
)

//nolint:gochecknoinits // That's how this is supposed to be done.
func init() {
	register.Plugin("arangolint", New)
}

// Arangolint is a golangci plugin.
type Arangolint struct{}

// New returns an arangolint linter as a golangci plugin.
//
//nolint:ireturn
func New(_ any) (register.LinterPlugin, error) {
	return &Arangolint{}, nil
}

// BuildAnalyzers implements register.LinterPlugin.
func (f *Arangolint) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.NewAnalyzer(),
	}, nil
}

// GetLoadMode implements register.LinterPlugin.
func (f *Arangolint) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
