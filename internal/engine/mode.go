package engine

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Mode int

const (
	SimpleMode Mode = iota
	ProMode
)

type Flags struct {
	Simple bool
	Pro    bool
}

// DetectMode decides Pro vs Simple from config.yaml when present.
// A file named config.yaml alone is not enough (many tools use that name in the
// user profile or project root); we require a recognizable Aeroform shape:
// mode: pro|simple, or cloud: aws|azure|gcp.
func DetectMode(cfgPath string, flags Flags) Mode {
	if flags.Simple {
		return SimpleMode
	}
	if flags.Pro {
		return ProMode
	}
	if _, err := os.Stat(cfgPath); err != nil {
		return SimpleMode
	}

	v := viper.New()
	v.SetConfigFile(cfgPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return SimpleMode
	}

	switch strings.ToLower(strings.TrimSpace(v.GetString("mode"))) {
	case "simple":
		return SimpleMode
	case "pro":
		return ProMode
	}

	switch strings.ToLower(strings.TrimSpace(v.GetString("cloud"))) {
	case "aws", "azure", "gcp":
		return ProMode
	}

	return SimpleMode
}
