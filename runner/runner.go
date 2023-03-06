package runner

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sv-tools/gochecker/config"
	"github.com/sv-tools/gochecker/output"
)

func Intercept() {
	conf := config.ParseConfig()

	buf := runMultiChecker(conf.Args...)
	// exit if no issues
	if buf.Len() == 3 && bytes.Equal(bytes.TrimSpace(buf.Bytes()), []byte("{}")) {
		return
	}
	diag := output.ParseOutput(conf, buf)
	if len(*diag) == 0 {
		os.Exit(0)
	}
	switch conf.Output {
	case config.ConsoleOutput:
		output.PrintAsConsole(diag)
		os.Exit(3)
	case config.JSONOutput:
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		if err := e.Encode(diag); err != nil {
			log.Fatalf("json ouput failed: %+v", err)
		}
	case config.GithubOutput:
		output.PrintAsGithub(diag)
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
	cmd.Env = append(os.Environ(), InterceptModeEnv+"=on")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	var code = -1
	if cmd.ProcessState != nil {
		code = cmd.ProcessState.ExitCode()
	}
	if stderr.Len() > 0 {
		if _, err = stderr.WriteTo(os.Stderr); err != nil {
			log.Printf("writing to stderr failed: %+v", err)
			if code == 0 {
				code = 1
			}
		}
		os.Exit(code)
	}
	if err != nil {
		log.Fatalf("interception failed: %+v", err)
	}
	return &stdout
}
