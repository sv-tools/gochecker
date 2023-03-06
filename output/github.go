package output

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func PrintAsGithub(diag *Diagnostic) {
	// print console output for people
	os.Stdout.WriteString("::group::console format\n")
	PrintAsConsole(diag)
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
						buf.WriteString(strings.Replace(strings.TrimSuffix(f.Lines[line-1], "\n"), "\t", " ", pos))
						if pos != -1 {
							buf.WriteString("%0A")
							buf.Grow(pos)
							for i := 0; i < pos-1; i++ {
								buf.WriteRune(' ')
							}
							buf.WriteRune('^')
						}
					}
					for _, fix := range issue.SuggestedFixes {
						buf.WriteString("%0A")
						buf.WriteString("Suggested Fix:")
						if fix.Message != "" {
							buf.WriteRune(' ')
							buf.WriteString(fix.Message)
						}
						buf.WriteString("%0A```diff%0A")
						buf.WriteString(strings.ReplaceAll(fix.Diff, "\n", "%0A"))
						buf.WriteString("```")
					}

					buf.WriteRune('\n')
					buf.WriteTo(os.Stdout)
				}()
			}
		}
	}
	wg.Wait()
}
