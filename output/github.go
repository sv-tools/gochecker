package output

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/sv-tools/gochecker/config"
)

func PrintAsGithub(diag *Diagnostic) (ret bool) {
	// print console output for people
	if _, err := os.Stdout.WriteString("::group::console format\n"); err != nil {
		log.Printf("writing to stdout failed: %+v", err)
		os.Exit(1)
	}
	ret = PrintAsConsole(diag)
	if _, err := os.Stdout.WriteString("::endgroup::\n"); err != nil {
		log.Printf("writing to stdout failed: %+v", err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	for _, pkg := range *diag {
		for name, obj := range pkg {
			if obj.Error != "" {
				buf := bytes.Buffer{}
				buf.WriteString("::error:: ")
				buf.WriteString(name)
				buf.WriteString(": ")
				buf.WriteString(obj.Error)
				buf.WriteRune('\n')
				if _, err := buf.WriteTo(os.Stdout); err != nil {
					log.Printf("writing to stdout failed: %+v", err)
					os.Exit(1)
				}
			}
			for _, issue := range obj.Issues {
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
					switch issue.SeverityLevel {
					case config.ErrorLevel:
						buf.WriteString("::error file=")
					case config.WarningLevel:
						buf.WriteString("::warning file=")
					case config.InfoLevel:
						buf.WriteString("::notice file=")
					}
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
					if _, err = buf.WriteTo(os.Stdout); err != nil {
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
