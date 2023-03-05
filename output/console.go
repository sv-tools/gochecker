package output

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

func PrintAsConsole(diag *Diagnostic) {
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
