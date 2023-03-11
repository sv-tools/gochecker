package output

import (
	"log"
	"os"
	"regexp"

	"github.com/sv-tools/gochecker/config"
)

var nolintRE = regexp.MustCompile(`//\s*nolint`)

func Modify(conf *config.Config, diag *Diagnostic) {
	toDeletePkg := make([]string, 0, len(*diag))
	for pkgName, pkg := range *diag {
		toDeleteAnalyzer := make([]string, 0, len(pkg))
		for analyzerName, obj := range pkg {
			if obj.Error != "" {
				continue
			}
			tmp := make([]*Issue, 0, len(obj.Issues))
			for _, issue := range obj.Issues {
				switch {
				case conf.Fix && len(issue.SuggestedFixes) != 0: // remove all issues with suggested fixes, because they are already applied
				case isNolint(issue): // remove issues with nolint comment
				case isExcluded(conf.Exclude, pkgName, analyzerName, issue):
				default:
					setSeverityLevel(conf.Severity, pkgName, analyzerName, issue)
					tmp = append(tmp, issue)
				}
			}
			if len(tmp) == 0 {
				toDeleteAnalyzer = append(toDeleteAnalyzer, analyzerName)
			} else {
				pkg[analyzerName].Issues = tmp
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

func isExcluded(rules []*config.Rule, pkg, analyzer string, issue *Issue) bool {
	for _, rule := range rules {
		if matchRule(rule, pkg, analyzer, issue) {
			return true
		}
	}
	return false
}

func matchRule(rule *config.Rule, pkg, analyzer string, issue *Issue) bool {
	if rule.Analyzer == "" && rule.PackageRE == nil && rule.PathRE == nil && rule.MessageRE == nil {
		return false
	}
	if rule.Analyzer != "" && analyzer != rule.Analyzer {
		return false
	}
	if rule.PackageRE != nil && !rule.PackageRE.MatchString(pkg) {
		return false
	}
	if rule.PathRE != nil && !rule.PathRE.MatchString(issue.PosN) {
		return false
	}
	if rule.MessageRE != nil && !rule.MessageRE.MatchString(issue.Message) {
		return false
	}
	return true
}

func setSeverityLevel(sevRules []*config.SeverityRule, pkg, analyzer string, issue *Issue) {
	issue.SeverityLevel = config.ErrorLevel
	for _, sev := range sevRules {
		for _, rule := range sev.Rules {
			if matchRule(rule, pkg, analyzer, issue) {
				issue.SeverityLevel = sev.Level
				return
			}
		}
	}
}
