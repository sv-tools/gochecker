package main

import (
	"log"
	"os"

	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/sv-tools/gochecker/analyzers"
	"github.com/sv-tools/gochecker/config"
)

const (
	progname         = "gochecker"
	interceptModeEnv = "330311d0-6a15-420d-8a66-a35b3e0f9c40"
)

func main() {
	os.Args[0] = progname
	log.SetFlags(0)
	log.SetPrefix(progname + ": ")

	if os.Getenv(interceptModeEnv) != "" {
		// pass directly to multichecker
		multichecker.Main(analyzers.Analyzers...)
	}

	// check for any sub-commands
	commands()

	// intercept the output of the multichecker and do the job
	intercept()
}

func commands() {
	if len(os.Args) == 1 {
		return
	}
	switch os.Args[1] {
	case "help":
		multichecker.Main(analyzers.Analyzers...)
	case "generate-config":
		config.GenerateConfig()
	}
}
