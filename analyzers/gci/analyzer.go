// Package gci integrates the [gci](https://github.com/daixiang0/gci) linter
// wrapping gci.Analyzer in order to provide the Suggested Fixes
// A lot of code has been copied from the gci/pkg/analyzer package.
package gci

import (
	"flag"
	"fmt"
	"go/token"
	"strconv"
	"strings"

	"github.com/daixiang0/gci/pkg/analyzer"
	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/io"
	"golang.org/x/tools/go/analysis"

	"github.com/sv-tools/gochecker/analyzers/utils"
)

var Analyzer = analyzer.Analyzer

func init() {
	Analyzer.Run = runAnalysis
}

func runAnalysis(pass *analysis.Pass) (any, error) {
	var fileReferences []*token.File
	// extract file references for all files in the analyzer pass
	for _, pkgFile := range pass.Files {
		fileForPos := pass.Fset.File(pkgFile.Package)
		if fileForPos != nil {
			fileReferences = append(fileReferences, fileForPos)
		}
	}
	expectedNumFiles := len(pass.Files)
	foundNumFiles := len(fileReferences)
	if expectedNumFiles != foundNumFiles {
		return nil, fmt.Errorf("expected %d files in Analyzer input, found %d", expectedNumFiles, foundNumFiles)
	}
	// read configuration options
	gciCfg, err := parseGciConfiguration(Analyzer.Flags)
	if err != nil {
		return nil, err
	}

	for _, file := range fileReferences {
		filePath := file.Name()
		unmodifiedFile, formattedFile, err := gci.LoadFormatGoFile(io.File{FilePath: filePath}, *gciCfg)
		if err != nil {
			return nil, err
		}
		fixes, err := utils.GetSuggestedFixesFromDiff(file, unmodifiedFile, formattedFile)
		if err != nil {
			return nil, err
		}
		if len(fixes) == 0 {
			// no difference
			continue
		}
		pass.Report(analysis.Diagnostic{
			Pos:            fixes[0].TextEdits[0].Pos,
			Message:        fmt.Sprintf("fix by `%s %s`", generateCmdLine(*gciCfg), filePath),
			SuggestedFixes: fixes,
		})
	}
	return nil, nil
}

func parseGciConfiguration(fs flag.FlagSet) (*config.Config, error) {
	var (
		noInlineComments     bool
		noPrefixComments     bool
		skipGenerated        bool
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
	if s := fs.Lookup(analyzer.SkipGeneratedFlag).Value.String(); s != "" {
		skipGenerated, err = strconv.ParseBool(s)
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
		SkipGenerated:    skipGenerated,
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

	if cfg.BoolConfig.SkipGenerated {
		result += " --skip-generated "
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
