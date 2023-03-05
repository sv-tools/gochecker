package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis/multichecker"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/gochecker/analyzers"
)

const (
	ConsoleOutput = "console"
	JSONOutput    = "json"
	GithubOutput  = "github"
)

var oneOfOutputFormats = strings.Join([]string{ConsoleOutput, JSONOutput, GithubOutput}, ", ")

type Config struct {
	Analyzers  map[string]map[string]string `json:"analyzers" yaml:"analyzers"`
	Debug      string                       `json:"debug" yaml:"debug"`
	CPUProfile string                       `json:"cpuprofile" yaml:"cpuprofile"`
	MemProfile string                       `json:"memprofile" yaml:"memprofile"`
	Trace      string                       `json:"trace" yaml:"trace"`
	Output     string                       `json:"output" yaml:"output"`
	Args       []string                     `json:"-" yaml:"-"`
	Test       bool                         `json:"test" yaml:"test"`
	Fix        bool                         `json:"fix" yaml:"fix"`
}

func parseConfig() *Config {
	var (
		configPath string
		config     Config
	)
	fs := flag.NewFlagSet(progname, flag.ContinueOnError)
	fs.SetOutput(&bytes.Buffer{}) // mute any prints
	var govetFlag bool
	for _, set := range []*flag.FlagSet{flag.CommandLine, fs} {
		set.StringVar(&configPath, "config", "", "A path to a config file in json or yaml format.")
		set.StringVar(&config.Output, "output", "", "Output format, one of: "+oneOfOutputFormats)
		set.BoolVar(&govetFlag, analyzers.GoVetName, false, analyzers.GoVetDoc)
	}
	// default flags multichecker flags
	fs.StringVar(&config.Debug, "debug", "", "")
	fs.StringVar(&config.CPUProfile, "cpuprofile", "", "")
	fs.StringVar(&config.MemProfile, "memprofile", "", "")
	fs.StringVar(&config.Trace, "trace", "", "")
	fs.BoolVar(&config.Test, "test", true, "")
	fs.BoolVar(&config.Fix, "fix", false, "")
	var jsonFlag bool
	fs.BoolVar(&jsonFlag, "json", false, "")
	// analyzer's flags
	for _, analyzer := range analyzers.Analyzers {
		fs.Bool(analyzer.Name, false, "")
		analyzer.Flags.VisitAll(func(f *flag.Flag) {
			fs.Var(f.Value, strings.Join([]string{analyzer.Name, f.Name}, "."), "")
		})
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		// let the multichecker report about any parser errors
		multichecker.Main(analyzers.Analyzers...)
	}
	if configPath != "" {
		f, err := os.Open(configPath)
		if err != nil {
			log.Fatal(err)
		}
		d := yaml.NewDecoder(f)
		d.KnownFields(true)
		if err := d.Decode(&config); err != nil {
			log.Fatal(err)
		}
	}
	if jsonFlag {
		config.Output = "json"
	}
	config.Output = strings.ToLower(config.Output)
	switch config.Output {
	case "":
		config.Output = ConsoleOutput
	case ConsoleOutput, JSONOutput, GithubOutput:
	default:
		log.Fatal("output must be one of: " + oneOfOutputFormats)
	}

	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "config", "output":
			return
		case "debug", "cpuprofile", "memprofile", "trace", "test", "fix", "json":
			return
		case analyzers.GoVetName:
			return
		}
		parts := strings.SplitN(f.Name, ".", 2)
		name := parts[0]
		flags, ok := config.Analyzers[name]
		if !ok {
			flags = make(map[string]string)
		}
		if len(parts) == 2 {
			flags[parts[1]] = f.Value.String()
		}
		config.Analyzers[name] = flags
	})
	if _, ok := config.Analyzers[analyzers.GoVetName]; ok || govetFlag {
		// Add all go vet linters to config
		delete(config.Analyzers, analyzers.GoVetName)
		for _, analyzer := range analyzers.GoVet {
			if _, ok := config.Analyzers[analyzer.Name]; !ok {
				config.Analyzers[analyzer.Name] = make(map[string]string)
			}
		}
	}

	args := []string{"-json"}
	if config.Test {
		args = append(args, "-test")
	}
	if config.Fix {
		args = append(args, "-fix")
	}
	if config.Debug != "" {
		args = append(args, "-debug", config.Debug)
	}
	if config.CPUProfile != "" {
		args = append(args, "-cpuprofile", config.CPUProfile)
	}
	if config.MemProfile != "" {
		args = append(args, "-memprofile", config.MemProfile)
	}
	if config.Trace != "" {
		args = append(args, "-trace", config.Trace)
	}
	for name, flags := range config.Analyzers {
		args = append(args, "-"+name)
		for fname, value := range flags {
			if value != "" {
				switch strings.ToLower(value) {
				case "false":
					continue
				case "true":
					args = append(args, fmt.Sprintf("-%s.%s", name, fname))
				default:
					args = append(args, fmt.Sprintf("-%s.%s", name, fname), value)
				}
			}
		}
	}
	config.Args = append(args, fs.Args()...)
	return &config
}

func generateConfig() {
	var config Config
	config.Analyzers = make(map[string]map[string]string)
	for _, analyzer := range analyzers.Analyzers {
		flags := make(map[string]string)
		analyzer.Flags.VisitAll(func(f *flag.Flag) {
			flags[f.Name] = f.DefValue
		})
		config.Analyzers[analyzer.Name] = flags
	}
	config.Analyzers[analyzers.GoVetName] = make(map[string]string)
	if err := yaml.NewEncoder(os.Stdout).Encode(config); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
