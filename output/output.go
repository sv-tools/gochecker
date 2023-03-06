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

func ParseOutput(conf *config.Config, data *bytes.Buffer) *Diagnostic {
	out := make(Diagnostic)
	d := json.NewDecoder(data)
	d.DisallowUnknownFields()
	if err := d.Decode(&out); err != nil {
		log.Fatalf("unmarshaling failed: %+v", err)
	}
	Exclude(conf, &out)
	return &out
}
