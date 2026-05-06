package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envlens/internal/auditor"
	"github.com/yourorg/envlens/internal/differ"
	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/reporter"
	"github.com/yourorg/envlens/internal/validator"
)

func main() {
	var (
		baseFile   = flag.String("base", "", "Base env file (e.g. .env.staging)")
		targetFile = flag.String("target", "", "Target env file (e.g. .env.production)")
		rulesFile  = flag.String("rules", "", "Optional validation rules JSON file")
		format     = flag.String("format", "text", "Output format: text or json")
		auditOnly  = flag.Bool("audit", false	, "Run audit only on target file")
	)
	flag.Parse()

	if *targetFile == "" {
		fmt.Fprintln(os.Stderr, "error: --target is required")
		flag.Usage()
		os.Exit(1)
	}

	targetEnv, err := parser.ParseFile(*targetFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing target file: %v\n", err)
		os.Exit(1)
	}

	var diffs []differ.DiffEntry
	if *baseFile != "" {
		baseEnv, err := parser.ParseFile(*baseFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing base file: %v\n", err)
			os.Exit(1)
		}
		diffs = differ.Diff(baseEnv, targetEnv)
	}

	var validationIssues []validator.Issue
	if *rulesFile != "" {
		rules, err := validator.LoadRules(*rulesFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading rules: %v\n", err)
			os.Exit(1)
		}
		validationIssues = validator.Validate(targetEnv, rules)
	}

	var auditIssues []auditor.Issue
	if *auditOnly || *baseFile == "" {
		auditIssues = auditor.Audit(targetEnv)
	}

	err = reporter.Write(os.Stdout, *format, diffs, validationIssues, auditIssues)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}

	if len(validationIssues) > 0 {
		os.Exit(2)
	}
}
