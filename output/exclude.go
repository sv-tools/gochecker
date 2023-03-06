package output

import (
	"log"
	"os"
	"regexp"

	"github.com/sv-tools/gochecker/config"
)

var nolintRE = regexp.MustCompile(`//\s*nolint`)

func Exclude(conf *config.Config, diag *Diagnostic) {
	toDeletePkg := make([]string, 0, len(*diag))
	for pkgName, pkg := range *diag {
		toDeleteAnalyzer := make([]string, 0, len(pkg))
		for analyzerName, issues := range pkg {
			tmp := make([]*Issue, 0, len(issues))
			for _, issue := range issues {
				switch {
				case conf.Fix && len(issue.SuggestedFixes) != 0: // remove all issues with suggested fixes, because they are already applied
				case isNolint(issue): // remove issues with nolint comment
				case isExclude(conf, pkgName, analyzerName, issue):
				default:
					tmp = append(tmp, issue)
				}
			}
			if len(tmp) == 0 {
				toDeleteAnalyzer = append(toDeleteAnalyzer, analyzerName)
			} else {
				pkg[analyzerName] = tmp
			}
		}
		for _, name := range toDeleteAnalyzer {
			delete(pkg, name)
		}
		if len(pkg) == 0 {
			toDeletePkg = append(toDeletePkg, pkgName)
		}
	}
	for _, name := range toDeletePkg {
		delete(*diag, name)
	}
}

func isNolint(issue *Issue) bool {
	filename, line, _ := parsePosN(issue.PosN)
	f, err := getFile(filename)
	if err != nil {
		log.Printf("reading file %q failed: %+v", filename, err)
		os.Exit(1)
	}
	return line != -1 && line < len(f.Lines) && nolintRE.MatchString(f.Lines[line-1])
}

func isExclude(conf *config.Config, pkg, analyzer string, issue *Issue) bool {
	for _, rule := range conf.Exclude {
		if rule.Analyzer != "" && analyzer != rule.Analyzer {
			continue
		}
		if rule.PackageRE != nil && !rule.PackageRE.MatchString(pkg) {
			continue
		}
		if rule.PathRE != nil && !rule.PathRE.MatchString(issue.PosN) {
			continue
		}
		if rule.MessageRE != nil && !rule.MessageRE.MatchString(issue.Message) {
			continue
		}
		return true
	}
	return false
}
