package config

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	critic "github.com/go-critic/go-critic/checkers/analyzer"
	"golang.org/x/tools/go/analysis/multichecker"
	"gopkg.in/yaml.v3"

	"github.com/sv-tools/gochecker/analyzers"
)

const (
	ConsoleOutput = "console"
	JSONOutput    = "json"
	GithubOutput  = "github"

	ErrorLevel   = "error"
	WarningLevel = "warning"
	InfoLevel    = "info"
)

var oneOfOutputFormats = strings.Join([]string{ConsoleOutput, JSONOutput, GithubOutput}, ", ")

type Config struct {
	Analyzers  map[string]map[string]string `json:"analyzers" yaml:"analyzers"`
	Module     string                       `json:"module" yaml:"module"`
	Debug      string                       `json:"debug" yaml:"debug"`
	CPUProfile string                       `json:"cpuprofile" yaml:"cpuprofile"`
	MemProfile string                       `json:"memprofile" yaml:"memprofile"`
	Trace      string                       `json:"trace" yaml:"trace"`
	Output     string                       `json:"output" yaml:"output"`
	GoVersion  string                       `json:"go_version" yaml:"go_version"`
	Args       []string                     `json:"-" yaml:"-"`
	Severity   []*SeverityRule              `json:"severity" yaml:"severity"`
	Exclude    []*Rule                      `json:"exclude" yaml:"exclude"`
	Test       bool                         `json:"test" yaml:"test"`
	Fix        bool                         `json:"fix" yaml:"fix"`
}

type Rule struct {
	PackageRE *regexp.Regexp `json:"-" yaml:"-"`
	PathRE    *regexp.Regexp `json:"-" yaml:"-"`
	MessageRE *regexp.Regexp `json:"-" yaml:"-"`
	Package   string         `json:"package" yaml:"package"`
	Analyzer  string         `json:"analyzer" yaml:"analyzer"`
	Path      string         `json:"path" yaml:"path"`
	Message   string         `json:"message" yaml:"message"`
	Severity  string         `json:"severity" yaml:"severity"`
	GitRef    string         `json:"git_ref" yaml:"git_ref"`
}

type SeverityRule struct {
	Level string  `json:"level" yaml:"level"`
	Rules []*Rule `json:"rules" yaml:"rules"`
}

func compileRules(exclude []*Rule) error {
	var err error
	for _, rule := range exclude {
		if rule.Package != "" {
			rule.PackageRE, err = regexp.Compile(rule.Package)
			if err != nil {
				return err
			}
		}
		if rule.Path != "" {
			rule.PathRE, err = regexp.Compile(rule.Path)
			if err != nil {
				return err
			}
		}
		if rule.Message != "" {
			rule.MessageRE, err = regexp.Compile(rule.Message)
			if err != nil {
				return err
			}
		}
		rule.Severity = strings.ToLower(rule.Severity)
		switch rule.Severity {
		case "":
		case ErrorLevel, WarningLevel, InfoLevel:
		default:
			return fmt.Errorf("wrong severity level %q, must be one of: %s", rule.Severity, strings.Join([]string{ErrorLevel, WarningLevel, InfoLevel}, ", "))
		}
	}
	return nil
}

func ParseConfig() *Config {
	var (
		configPath string
		config     = Config{Analyzers: map[string]map[string]string{}}
	)
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(&bytes.Buffer{}) // mute any prints
	for _, set := range []*flag.FlagSet{flag.CommandLine, fs} {
		set.StringVar(&configPath, "config", "", "A path to a config file in json or yaml format.")
		set.StringVar(&config.Output, "output", "", "Output format, one of: "+oneOfOutputFormats)
	}
	// default multichecker flags
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
		if err = d.Decode(&config); err != nil {
			log.Fatal(err)
		}
		if err := compileRules(config.Exclude); err != nil {
			log.Fatal(err)
		}
		for _, sev := range config.Severity {
			sev.Level = strings.ToLower(sev.Level)
			switch sev.Level {
			case "":
				sev.Level = ErrorLevel
			case ErrorLevel, WarningLevel, InfoLevel:
			default:
				log.Fatalf("wrong severity level %q, must be one of: %s", sev.Level, strings.Join([]string{ErrorLevel, WarningLevel, InfoLevel}, ", "))
			}
			if err := compileRules(sev.Rules); err != nil {
				log.Fatal(err)
			}
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
	if v, ok := config.Analyzers[analyzers.GoVetName]; ok {
		excludes := make(map[string]struct{})
		govetAnalyzers := analyzers.GoVet
		if s, ok := v[analyzers.GoVetExtraName]; ok {
			t, err := strconv.ParseBool(s)
			if err != nil {
				log.Fatalf("unable to parse govet.extra: %+v", err)
			}
			if t {
				govetAnalyzers = append(govetAnalyzers, analyzers.GoVetExtra...)
			}
		}
		if s := v[analyzers.GoVetExcludeName]; s != "" {
			vetAnalyzers := make(map[string]struct{}, len(govetAnalyzers))
			for _, a := range govetAnalyzers {
				vetAnalyzers[a.Name] = struct{}{}
			}
			for _, exc := range strings.Split(s, ",") {
				name := strings.TrimSpace(exc)
				if _, ok := vetAnalyzers[name]; !ok {
					log.Fatalf("analyzer %q is not a part of go vet passes", name)
				}
				excludes[name] = struct{}{}
			}
		}
		// Add all go vet linters to config
		delete(config.Analyzers, analyzers.GoVetName)
		for _, analyzer := range govetAnalyzers {
			name := analyzer.Name
			if _, ok := excludes[name]; ok {
				delete(config.Analyzers, name)
				continue
			}
			if _, ok = config.Analyzers[name]; !ok {
				config.Analyzers[name] = make(map[string]string)
			}
		}
	}

	if err := ApplyModInfo(&config); err != nil {
		log.Fatal("Reading info about go.mo failed: #+v", err)
	}

	// preparing the args

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

func GenerateConfig() {
	var config Config
	config.Analyzers = make(map[string]map[string]string)
	for _, analyzer := range analyzers.Analyzers {
		flags := make(map[string]string)
		analyzer.Flags.VisitAll(func(f *flag.Flag) {
			flags[f.Name] = f.Value.String()
		})
		config.Analyzers[analyzer.Name] = flags
	}
	config.Exclude = []*Rule{
		{
			Analyzer: "",
			Path:     "",
			Package:  "",
			Message:  "",
		},
	}
	config.Severity = []*SeverityRule{
		{
			Level: "error",
			Rules: []*Rule{
				{
					Analyzer: "",
					Path:     "",
					Package:  "",
					Message:  "",
				},
			},
		},
	}
	// set the default value in example config to 8 to avoid the difference between local and GitHub checks
	config.Analyzers[critic.Analyzer.Name]["concurrency"] = "8"

	if err := yaml.NewEncoder(os.Stdout).Encode(config); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
