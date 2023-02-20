package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Analyzers  map[string]map[string]string `json:"analyzers" yaml:"analyzers"`
	Debug      string                       `json:"debug" yaml:"debug"`
	CPUProfile string                       `json:"cpuprofile" yaml:"cpuprofile"`
	MemProfile string                       `json:"memprofile" yaml:"memprofile"`
	Trace      string                       `json:"trace" yaml:"trace"`
	Test       bool                         `json:"test" yaml:"test"`
	Fix        bool                         `json:"fix" yaml:"fix"`
}

const progname = "gochecker"

func parseConfig() {
	os.Args[0] = progname
	log.SetFlags(0)
	log.SetPrefix(progname + ": ")

	var (
		configPath string
		config     Config
	)
	configFlagSet := flag.NewFlagSet(progname, flag.ContinueOnError)
	configFlagSet.Usage = func() {} // preventing the help output
	for _, set := range []*flag.FlagSet{flag.CommandLine, configFlagSet} {
		set.StringVar(&configPath, "config", "", "A path to a config file in json or yaml format.")
	}
	// default flags multichecker flags
	configFlagSet.StringVar(&config.Debug, "debug", "", "")
	configFlagSet.StringVar(&config.CPUProfile, "cpuprofile", "", "")
	configFlagSet.StringVar(&config.MemProfile, "memprofile", "", "")
	configFlagSet.StringVar(&config.Trace, "trace", "", "")
	configFlagSet.BoolVar(&config.Test, "test", true, "")
	configFlagSet.BoolVar(&config.Fix, "fix", false, "")
	// analyzer flags
	for _, analyzer := range analyzers {
		configFlagSet.Bool(analyzer.Name, false, "")
		analyzer.Flags.VisitAll(func(f *flag.Flag) {
			configFlagSet.Var(f.Value, strings.Join([]string{analyzer.Name, f.Name}, "."), "")
		})
	}
	if err := configFlagSet.Parse(os.Args[1:]); err != nil {
		return // let the multichecker report about any errors
	}

	if largs := configFlagSet.Args(); len(largs) > 0 {
		switch largs[0] {
		case "generate-config":
			generateConfig()
			os.Exit(0)
		}
	}

	if configPath == "" {
		return // no config provided
	}
	f, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	d := yaml.NewDecoder(f)
	d.KnownFields(true)
	if err := d.Decode(&config); err != nil {
		log.Fatal(err)
	}
	configFlagSet.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "config", "debug", "cpuprofile", "memprofile", "trace", "test", "fix":
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
	args := []string{progname}
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
	args = append(args, configFlagSet.Args()...)
	os.Args = args
}

func generateConfig() {
	var config Config
	config.Analyzers = make(map[string]map[string]string)
	for _, analyzer := range analyzers {
		flags := make(map[string]string)
		analyzer.Flags.VisitAll(func(f *flag.Flag) {
			flags[f.Name] = f.DefValue
		})
		config.Analyzers[analyzer.Name] = flags
	}
	if err := yaml.NewEncoder(os.Stdout).Encode(config); err != nil {
		log.Fatal(err)
	}
}
