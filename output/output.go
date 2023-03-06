package output

import (
	"bytes"
	"encoding/json"
	"go/token"
	"log"

	"github.com/sv-tools/gochecker/config"
)

type (
	// Diagnostic is a structure to hold the json output from multichecker
	//
	// JSON example:
	// ```json
	//	{
	//	 "<package>": {
	//	   "<analyzer>": {
	//	     "posn": "</path/to/file.go>:<line num>:<column num>",
	//	     "message": "<message>",
	//	     "suggested_fixes": [
	//	       {
	//	         "message": "",
	//	         "edits": [
	//	           {
	//	             "filename": "</path/to/file.go>",
	//	             "start": 865,
	//	             "end": 865,
	//	             "new": "<new text>"
	//	           }
	//	         ]
	//	       }
	//	     ]
	//	   }
	//	 }
	//	}
	// ```
	Diagnostic map[string]map[string][]*Issue
	Issue      struct {
		Message        string `json:"message"`
		Category       string `json:"category,omitempty"`
		PosN           string `json:"posn"`
		SuggestedFixes []*Fix `json:"suggested_fixes,omitempty"`
	}
	Fix struct {
		Message string  `json:"message,omitempty"`
		Diff    string  `json:"-"` // system field to contain calculated diff
		Edits   []*Edit `json:"edits"`
	}
	Edit struct {
		Filename string    `json:"filename"`
		New      string    `json:"new"`
		Start    token.Pos `json:"start"`
		End      token.Pos `json:"end,omitempty"`
	}
)

func ParseOutput(data *bytes.Buffer) *Diagnostic {
	out := make(Diagnostic)
	d := json.NewDecoder(data)
	d.DisallowUnknownFields()
	if err := d.Decode(&out); err != nil {
		log.Fatalf("unmarshaling failed: %+v", err)
	}
	return &out
}

func Exclude(conf *config.Config, diag *Diagnostic) {
	if conf.Fix {
		// remove all issues with suggested fixes, because they are already applied
		toDeletePkg := make([]string, 0, len(*diag))
		for pkgName, pkg := range *diag {
			toDeleteAnalyzer := make([]string, 0, len(pkg))
			for analyzerName, issues := range pkg {
				tmp := make([]*Issue, 0, len(issues))
				for _, issue := range issues {
					if len(issue.SuggestedFixes) == 0 {
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
}
