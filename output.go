package main

import (
	"bytes"
	"encoding/json"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/sv-tools/gochecker/config"
)

func intercept() {
	conf := config.ParseConfig()

	buf := runMultiChecker(conf.Args...)
	// exit if no issues
	if buf.Len() == 3 && bytes.Equal(bytes.TrimSpace(buf.Bytes()), []byte("{}")) {
		return
	}
	diag := parseOutput(buf)
	exclude(conf, diag)
	if len(*diag) == 0 {
		os.Exit(0)
	}
	switch conf.Output {
	case config.ConsoleOutput:
		printAsText(diag)
		os.Exit(3)
	case config.JSONOutput:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		if err := e.Encode(diag); err != nil {
			log.Fatalf("json ouput failed: %+v", err)
		}
	case config.GithubOutput:
		printAsGithub(diag)
		os.Exit(3)
	}
}

func runMultiChecker(args ...string) *bytes.Buffer {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	prog, err := os.Executable()
	if err != nil {
		log.Fatalf("getting execuable failed: %+v", err)
	}
	prog, err = filepath.EvalSymlinks(prog)
	if err != nil {
		log.Fatalf("evaluating symlink failed: %+v", err)
	}

	cmd := exec.Command(prog, args...)
	cmd.Env = append(os.Environ(), interceptModeEnv+"=on")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	var code = -1
	if cmd.ProcessState != nil {
		code = cmd.ProcessState.ExitCode()
	}
	if stderr.Len() > 0 {
		stderr.WriteTo(os.Stderr)
		os.Exit(code)
	}
	if err != nil {
		log.Fatalf("interception failed: %+v", err)
	}
	return &stdout
}

type (
	Diagnostic map[string]map[string][]Issue
	Issue      struct {
		Message        string `json:"message"`
		Category       string `json:"category,omitempty"`
		PosN           string `json:"posn"`
		SuggestedFixes []Fix  `json:"suggested_fixes,omitempty"`
	}
	Fix struct {
		Message string `json:"message,omitempty"`
		Edits   []Edit `json:"edits"`
	}
	Edit struct {
		Filename string    `json:"filename"`
		New      string    `json:"new"`
		Start    token.Pos `json:"start"`
		End      token.Pos `json:"end,omitempty"`
	}
)

func parseOutput(data *bytes.Buffer) *Diagnostic {
	out := make(Diagnostic)
	d := json.NewDecoder(data)
	d.DisallowUnknownFields()
	if err := d.Decode(&out); err != nil {
		log.Fatalf("unmarshaling failed: %+v", err)
	}
	return &out
}

func exclude(conf *config.Config, diag *Diagnostic) {
	if conf.Fix {
		// remove all issues with suggested fixes, because they are already applied
		toDeletePkg := make([]string, 0, len(*diag))
		for pkgName, pkg := range *diag {
			toDeleteAnalyzer := make([]string, 0, len(pkg))
			for analyzerName, issues := range pkg {
				tmp := make([]Issue, 0, len(issues))
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

type CachedFile struct {
	Filename string
	Data     []byte
	Lines    []string
}

var cache = sync.Map{}

func getFile(filename string) (*CachedFile, error) {
	obj, ok := cache.Load(filename)
	if !ok {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		rel, err := filepath.Rel(wd, filename)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(data), "\n")
		f := &CachedFile{
			Filename: rel,
			Data:     data,
			Lines:    lines,
		}
		obj, _ = cache.LoadOrStore(filename, f)
	}
	return obj.(*CachedFile), nil
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

func printAsText(diag *Diagnostic) {
	wg := sync.WaitGroup{}
	for _, pkg := range *diag {
		for name, issues := range pkg {
			for _, issue := range issues {
				wg.Add(1)
				name := name
				issue := issue
				go func() {
					defer wg.Done()
					filename, line, pos := parsePosN(issue.PosN)

					f, err := getFile(filename)
					if err != nil {
						log.Printf("reading file %q failed: %+v", filename, err)
						return
					}
					buf := bytes.Buffer{}
					buf.WriteString(f.Filename)
					if line != -1 {
						buf.WriteString(":" + strconv.Itoa(line))
						if pos != -1 {
							buf.WriteString(":" + strconv.Itoa(pos))
						}
					}
					buf.WriteString(colorRed)
					if issue.Category != "" {
						buf.WriteString(": ")
						buf.WriteString(issue.Category)
					}
					if issue.Message != "" {
						buf.WriteString(": ")
						buf.WriteString(issue.Message)
					}
					buf.WriteString(colorReset)
					buf.WriteString(" (")
					buf.WriteString(name)
					buf.WriteRune(')')

					if line != -1 && line < len(f.Lines) {
						buf.WriteRune('\n')
						buf.WriteString(strings.Replace(f.Lines[line-1], "\t", " ", pos))
						if pos != -1 {
							buf.WriteRune('\n')
							buf.Grow(pos)
							for i := 0; i < pos-1; i++ {
								buf.WriteRune(' ')
							}
							buf.WriteString(colorYellow)
							buf.WriteRune('^')
							buf.WriteString(colorReset)
						}
					}

					buf.WriteRune('\n')
					buf.WriteTo(os.Stdout)
				}()
			}
		}
	}
	wg.Wait()
}

func parsePosN(posN string) (string, int, int) {
	var (
		filename string
		line     = -1
		pos      = -1
		err      error
	)
	parts := strings.Split(posN, ":")
	l := len(parts)
	switch {
	case l > 3:
		parts = []string{
			strings.Join(parts[:l-2], ":"),
			parts[l-3],
			parts[l-2],
		}
		fallthrough
	case l == 3:
		pos, err = strconv.Atoi(parts[2])
		if err != nil {
			log.Printf("converting the position failed: %+v", err)
			pos = -1
		}
		fallthrough
	case l == 2:
		line, err = strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("converting the line number failed: %+v", err)
			line = -1
		}
		fallthrough
	default:
		filename = parts[0]
	}
	return filename, line, pos
}

func printAsGithub(diag *Diagnostic) {
	// print console output for people
	os.Stdout.WriteString("::group::console format\n")
	printAsText(diag)
	os.Stdout.WriteString("::endgroup::\n")

	wg := sync.WaitGroup{}
	for _, pkg := range *diag {
		for name, issues := range pkg {
			for _, issue := range issues {
				wg.Add(1)
				name := name
				issue := issue
				go func() {
					defer wg.Done()
					filename, line, pos := parsePosN(issue.PosN)

					f, err := getFile(filename)
					if err != nil {
						log.Printf("reading file %q failed: %+v", filename, err)
						return
					}

					buf := bytes.Buffer{}
					buf.WriteString("::error file=")
					buf.WriteString(f.Filename)
					if line != -1 {
						buf.WriteString(",line=")
						buf.WriteString(strconv.Itoa(line))
						if pos != -1 {
							buf.WriteString(",col=")
							buf.WriteString(strconv.Itoa(pos))
						}
					}
					buf.WriteString("::")
					if issue.Category != "" {
						buf.WriteString(issue.Category)
						buf.WriteString(": ")
					}
					if issue.Message != "" {
						buf.WriteString(issue.Message)
					}
					buf.WriteString(" (")
					buf.WriteString(name)
					buf.WriteRune(')')
					if line != -1 && line < len(f.Lines) {
						buf.WriteString("%0A")
						buf.WriteString(strings.Replace(f.Lines[line-1], "\t", " ", pos))
						if pos != -1 {
							buf.WriteString("%0A")
							buf.Grow(pos)
							for i := 0; i < pos-1; i++ {
								buf.WriteRune(' ')
							}
							buf.WriteRune('^')
						}
					}

					buf.WriteRune('\n')
					buf.WriteTo(os.Stdout)
				}()
			}
		}
	}
	wg.Wait()
}
