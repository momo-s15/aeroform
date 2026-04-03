package engine

import "os"

type Mode int

const (
	SimpleMode Mode = iota
	ProMode
)

type Flags struct {
	Simple bool
	Pro    bool
}

func DetectMode(cfgPath string, flags Flags) Mode {
	if flags.Simple {
		return SimpleMode
	}
	if flags.Pro {
		return ProMode
	}
	if _, err := os.Stat(cfgPath); err == nil {
		return ProMode
	}
	return SimpleMode
}
