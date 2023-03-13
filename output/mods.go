package output

import (
	"bytes"
	"io"
	"log"
	"os"
	"regexp"
	"sort"

	"github.com/sv-tools/gochecker/config"
)

var nolintRE = regexp.MustCompile(`//\s*nolint`)

func Modify(conf *config.Config, diag *Diagnostic) {
	toFix := make(map[string][]*Edit)
	toDeletePkg := make([]string, 0, len(*diag))
	for pkgName, pkg := range *diag {
		toDeleteAnalyzer := make([]string, 0, len(pkg))
		for analyzerName, obj := range pkg {
			if obj.Error != "" {
				continue
			}
			tmp := make([]*Issue, 0, len(obj.Issues))
			for _, issue := range obj.Issues {
				setSeverityLevel(conf.Severity, pkgName, analyzerName, issue)
				switch {
				case isNolint(issue): // remove issues with nolint comment
				case isExcluded(conf.Exclude, pkgName, analyzerName, issue):
				case conf.Fix && len(issue.SuggestedFixes) > 0: // must be last in the order, so other rules are applied
					for _, fix := range issue.SuggestedFixes {
						for _, edit := range fix.Edits {
							toFix[edit.Filename] = append(toFix[edit.Filename], edit)
						}
					}
				default:
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
	if len(toFix) > 0 {
		ApplySuggestedFixes(toFix)
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
	if rule.Analyzer == "" && rule.PackageRE == nil && rule.PathRE == nil && rule.MessageRE == nil && rule.Severity == "" {
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
	if rule.Severity != "" && rule.Severity != issue.SeverityLevel {
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

func ApplySuggestedFixes(fixes map[string][]*Edit) {
	for filename, edits := range fixes {
		sort.Slice(edits, func(i, j int) bool {
			return edits[i].Start < edits[j].Start
		})
		f, err := getFile(filename)
		if err != nil {
			log.Fatalf("reading file %q failed: %+v", filename, err)
		}
		var end int64
		r := bytes.NewReader(f.Data)
		buf := bytes.Buffer{}
		for _, edit := range edits {
			l := int64(edit.Start) - end
			switch {
			case l < 0:
				log.Fatalf("overlapped change for file %q: %#v", filename, edit)
			case l > 0:
				b := make([]byte, l)
				if _, err := r.Read(b); err != nil {
					log.Fatalf("reading failed: %+v", err)
				}
				buf.Write(b)
			}
			buf.WriteString(edit.New)
			if end, err = r.Seek(int64(edit.End), io.SeekStart); err != nil {
				log.Fatalf("seeking failed: %+v", end)
			}
		}
		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatalf("reading leftovers failed: %#v", err)
		}
		if len(b) > 0 {
			buf.Write(b)
		}
		f.Data = buf.Bytes()
		if err := os.WriteFile(f.Filename, f.Data, 0o644); err != nil {
			log.Fatalf("writing to file %q failed: %#v", f.Filename, err)
		}
	}
}
