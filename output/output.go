package output

import (
	"bytes"
	"encoding/json"
	"go/token"
	"io"
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
	Diagnostic    map[string]map[string]*IssuesOrError
	IssuesOrError struct {
		Error  string
		Issues []*Issue
	}
	Issue struct {
		Message        string `json:"message"`
		Category       string `json:"category,omitempty"`
		PosN           string `json:"posn"`
		SeverityLevel  string `json:"severity_level"`
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

func (o *IssuesOrError) UnmarshalJSON(data []byte) error {
	var e struct {
		Error string `json:"error"`
	}
	r := bytes.NewReader(data)
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	if d.Decode(&e) == nil {
		o.Error = e.Error
		return nil
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}
	return d.Decode(&o.Issues)
}

func (o *IssuesOrError) MarshalJSON() ([]byte, error) {
	if o.Error != "" {
		e := struct {
			Error string `json:"error"`
		}{o.Error}
		return json.Marshal(e)
	}
	return json.Marshal(o.Issues)
}

func ParseOutput(conf *config.Config, data []byte) *Diagnostic {
	out := make(Diagnostic)
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&out); err != nil {
		log.Fatalf("unmarshaling failed \"%+v\" for response:\n%s", err, string(data))
	}
	Modify(conf, &out)
	return &out
}
