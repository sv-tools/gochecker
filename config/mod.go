package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// ModInfo contains the version of go and module name
type ModInfo struct {
	Path      string `json:"Path"`
	GoVersion string `json:"GoVersion"`
}

// GetModInfo returns the base info about module name and go version
// TODO: support path in case of monorepo
func GetModInfo() (*ModInfo, error) {
	cmd := exec.Command("go", "list", "-m", "-json")
	cmd.Env = os.Environ()
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("go mod failed '%w' for data: %s", err, string(data))
	}
	var mod ModInfo
	if err := json.Unmarshal(data, &mod); err != nil {
		return nil, fmt.Errorf("unmarshal data failed '%w' for data: %s", err, string(data))
	}

	return &mod, nil
}

// ApplyModInfo add the the mod info to all analyzers are in need of such info
func ApplyModInfo(conf *Config) error {
	if conf.Module == "" || conf.GoVersion == "" {
		mod, err := GetModInfo()
		if err != nil {
			return err
		}
		if conf.Module == "" {
			conf.Module = mod.Path
		}
		if conf.GoVersion == "" {
			conf.GoVersion = mod.GoVersion
		}
	}

	return nil
}
