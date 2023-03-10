package output

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/pmezard/go-difflib/difflib"

	"github.com/sv-tools/gochecker/config"
)

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
		lines := difflib.SplitLines(string(data))
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
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorPurple = "\033[35m"
)

// PrintAsConsole prinst the issues to console with colors and returns false in case only warning or info isues
func PrintAsConsole(diag *Diagnostic) (ret bool) {
	wg := sync.WaitGroup{}
	for _, pkg := range *diag {
		for name, obj := range pkg {
			if obj.Error != "" {
				buf := bytes.Buffer{}
				buf.WriteString(name)
				buf.WriteString(": ")
				buf.WriteString(colorRed)
				buf.WriteString(obj.Error)
				buf.WriteString(colorReset)
				buf.WriteRune('\n')
				if _, err := buf.WriteTo(os.Stdout); err != nil {
					log.Printf("writing to stdout failed: %+v", err)
					os.Exit(1)
				}
				continue
			}
			for _, issue := range obj.Issues {
				if issue.SeverityLevel == config.ErrorLevel {
					ret = true
				}
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
					switch issue.SeverityLevel {
					case config.ErrorLevel:
						buf.WriteString(colorRed)
						buf.WriteString(": ERR")
					case config.WarningLevel:
						buf.WriteString(colorYellow)
						buf.WriteString(": WRN")
					case config.InfoLevel:
						buf.WriteString(colorGreen)
						buf.WriteString(": INF")
					}
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
				FIXES:
					for _, fix := range issue.SuggestedFixes {
						buf.WriteString("Suggested Fix:")
						if fix.Message != "" {
							buf.WriteString(colorRed)
							buf.WriteRune(' ')
							buf.WriteString(fix.Message)
							buf.WriteString(colorReset)
						}
						buf.WriteRune('\n')
						reader := bytes.NewReader(f.Data)
						fixed := bytes.Buffer{}
						for _, edit := range fix.Edits {
							if edit.Filename != filename {
								// do not support modifications in multiple files for simplicity
								log.Printf("suggested fix for a file %q modifies another file %q: %#v", filename, edit.Filename, fix)
								break FIXES
							}
							var cur int64
							cur, err = reader.Seek(0, io.SeekCurrent)
							if err != nil {
								log.Printf("seeking on buffer for file %q failed: %+v", filename, err)
								break FIXES
							}
							if l := int64(edit.Start) - cur; l > 0 {
								b := make([]byte, l)
								if _, err = reader.Read(b); err != nil {
									log.Printf("reading from buffer for file %q failed: %+v", filename, err)
									break FIXES
								}
								fixed.Write(b)
							}
							fixed.WriteString(edit.New)
							end := edit.End
							if end < edit.Start {
								end = edit.Start
							}
							if _, err = reader.Seek(int64(end), io.SeekStart); err != nil {
								log.Printf("seeking on buffer for file %q failed: %+v", filename, err)
								break FIXES
							}
						}
						// read remaining data
						var data []byte
						data, err = io.ReadAll(reader)
						if err != nil {
							log.Printf("reading remaining data from buffer for file %q failed: %+v", filename, err)
							break
						}
						if len(data) > 0 {
							fixed.Write(data)
						}
						d := difflib.UnifiedDiff{
							A:       f.Lines,
							B:       difflib.SplitLines(fixed.String()),
							Context: 1,
						}
						var diff string
						diff, err = difflib.GetUnifiedDiffString(d)
						if err != nil {
							log.Printf("getting diff for file %q failed: %+v", filename, err)
							break FIXES
						}
						if diff == "" {
							continue
						}
						fix.Diff = diff
						lines := difflib.SplitLines(diff)
						for i := 0; i < len(lines); i++ {
							reset := false
							s := lines[i]
							switch {
							case strings.HasPrefix(s, "@@"):
								buf.WriteString(colorPurple)
								reset = true
							case strings.HasPrefix(s, "+"):
								reset = true
								buf.WriteString(colorGreen)
							case strings.HasPrefix(s, "-"):
								reset = true
								buf.WriteString(colorRed)
							}
							buf.WriteString(lines[i])
							if reset {
								buf.WriteString(colorReset)
							}
						}
					}
					if _, err := buf.WriteTo(os.Stdout); err != nil {
						log.Printf("writing to stdout failed: %+v", err)
						os.Exit(1)
					}
				}()
			}
		}
	}
	wg.Wait()
	return
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
