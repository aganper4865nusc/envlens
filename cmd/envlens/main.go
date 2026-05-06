package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envlens/internal/auditor"
	"github.com/yourorg/envlens/internal/differ"
	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/reporter"
	"github.com/yourorg/envlens/internal/resolver"
	"github.com/yourorg/envlens/internal/validator"
)

func main() {
	var (
		source        = flag.String("source", "", "source .env file (baseline)")
		target        = flag.String("target", "", "target .env file to compare/audit")
		format        = flag.String("format", "text", "output format: text|json")
		auditOnly     = flag.Bool("audit-only", false, "skip diff, only audit target")
		allowEnvOverride = flag.Bool("env-override", false, "let OS env vars override file values")
	)
	flag.Parse()

	if *target == "" {
		fmt.Fprintln(os.Stderr, "error: --target is required")
		os.Exit(1)
	}

	targetEnv, err := parser.ParseFile(*target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing target: %v\n", err)
		os.Exit(1)
	}

	// Resolve target env with optional OS override.
	resolutions, _ := resolver.Resolve(targetEnv, resolver.Options{
		AllowEnvOverride: *allowEnvOverride,
	})
	resolvedEnv := resolver.ToMap(resolutions)

	var diffs []differ.Entry
	var valIssues []validator.Issue

	if !*auditOnly && *source != "" {
		sourceEnv, err := parser.ParseFile(*source)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing source: %v\n", err)
			os.Exit(1)
		}
		diffs = differ.Diff(sourceEnv, resolvedEnv)

		// Validate that all source keys exist in target.
		required := make([]string, 0, len(sourceEnv))
		for k := range sourceEnv {
			required = append(required, k)
		}
		valIssues = validator.Validate(resolvedEnv, validator.Rules{Required: required})
	}

	findings := auditor.Audit(resolvedEnv)

	rep := reporter.Report{
		Diffs:       diffs,
		Validations: valIssues,
		Audits:      findings,
		Resolutions: resolutions,
	}

	if err := reporter.Write(os.Stdout, rep, *format); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}
}
