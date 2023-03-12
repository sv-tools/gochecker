// Package gci integrates the [gci](https://github.com/daixiang0/gci) linter
// wrapping gci.Analyzer in order to provide the Suggested Fixes
// A lot of code has been copied from the gci/pkg/analyzer package.
package gci

import (
	"flag"
	"fmt"
	"go/ast"
	"strconv"
	"strings"

	"github.com/daixiang0/gci/pkg/analyzer"
	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/io"
	"golang.org/x/tools/go/analysis"

	"github.com/sv-tools/gochecker/analyzers/skipgenerated"
	"github.com/sv-tools/gochecker/analyzers/utils"
)

var Analyzer = analyzer.Analyzer

func init() {
	Analyzer.Run = runAnalysis
	Analyzer.Requires = []*analysis.Analyzer{
		skipgenerated.Analyzer,
	}

	// removing the skipGenerated flag
	fs := flag.NewFlagSet(Analyzer.Flags.Name(), Analyzer.Flags.ErrorHandling())
	Analyzer.Flags.VisitAll(func(f *flag.Flag) {
		if f.Name == analyzer.SkipGeneratedFlag {
			return
		}
		fs.Var(f.Value, f.Name, f.Usage)
	})
	Analyzer.Flags = *fs
}

func runAnalysis(pass *analysis.Pass) (any, error) {
	files := pass.ResultOf[skipgenerated.Analyzer].([]*ast.File)
	if len(files) == 0 {
		return nil, nil
	}

	gciCfg, err := parseGciConfiguration(Analyzer.Flags)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileRef := pass.Fset.File(file.Pos())
		filePath := fileRef.Name()
		unmodifiedFile, formattedFile, err := gci.LoadFormatGoFile(io.File{FilePath: filePath}, *gciCfg)
		if err != nil {
			return nil, err
		}
		fix, err := utils.GetSuggestedFix(fileRef, unmodifiedFile, formattedFile)
		if err != nil {
			return nil, err
		}
		if fix == nil {
			// no difference
			continue
		}
		pass.Report(analysis.Diagnostic{
			Pos:            fix.TextEdits[0].Pos,
			Message:        fmt.Sprintf("fix by `%s %s`", generateCmdLine(*gciCfg), filePath),
			SuggestedFixes: []analysis.SuggestedFix{*fix},
		})
	}
	return nil, nil
}

func parseGciConfiguration(fs flag.FlagSet) (*config.Config, error) {
	var (
		noInlineComments     bool
		noPrefixComments     bool
		sectionsStr          string
		sectionSeparatorsStr string

		err error
	)

	if s := fs.Lookup(analyzer.NoInlineCommentsFlag).Value.String(); s != "" {
		noInlineComments, err = strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
	}
	if s := fs.Lookup(analyzer.NoPrefixCommentsFlag).Value.String(); s != "" {
		noPrefixComments, err = strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
	}
	sectionsStr = fs.Lookup(analyzer.SectionsFlag).Value.String()
	sectionSeparatorsStr = fs.Lookup(analyzer.SectionSeparatorsFlag).Value.String()

	fmtCfg := config.BoolConfig{
		NoInlineComments: noInlineComments,
		NoPrefixComments: noPrefixComments,
		Debug:            false,
		SkipGenerated:    false, // should be set to `false` to avoid unneeded scanning the file for the generated tokens
	}

	var sectionStrings []string
	if sectionsStr != "" {
		sectionStrings = strings.Split(sectionsStr, analyzer.SectionDelimiter)
	}

	var sectionSeparatorStrings []string
	if sectionSeparatorsStr != "" {
		sectionSeparatorStrings = strings.Split(sectionSeparatorsStr, analyzer.SectionDelimiter)
		fmt.Println(sectionSeparatorsStr)
	}
	return config.YamlConfig{Cfg: fmtCfg, SectionStrings: sectionStrings, SectionSeparatorStrings: sectionSeparatorStrings}.Parse()
}

func generateCmdLine(cfg config.Config) string {
	result := "gci write"

	if cfg.BoolConfig.NoInlineComments {
		result += " --NoInlineComments "
	}

	if cfg.BoolConfig.NoPrefixComments {
		result += " --NoPrefixComments "
	}

	if cfg.BoolConfig.CustomOrder {
		result += " --custom-order "
	}

	for _, s := range cfg.Sections.String() {
		result += fmt.Sprintf(" --Section \"%s\" ", s)
	}
	for _, s := range cfg.SectionSeparators.String() {
		result += fmt.Sprintf(" --SectionSeparator %s ", s)
	}
	return result
}
