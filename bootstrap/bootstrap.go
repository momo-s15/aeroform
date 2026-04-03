package bootstrap

import (
	"fmt"
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/providers"
)

type Report struct {
	Cloud string
	Steps []string
}

func Run(cfg config.Config) Report {
	provider := providers.ForCloud(cfg.Cloud)
	steps := provider.BootstrapSteps(cfg)
	if len(steps) == 0 {
		steps = []string{"No bootstrap steps defined yet"}
	}

	return Report{
		Cloud: provider.Name(),
		Steps: steps,
	}
}

func Render(report Report) string {
	var builder strings.Builder
	builder.WriteString("Aeroform bootstrap\n")
	builder.WriteString(fmt.Sprintf("cloud: %s\n", report.Cloud))
	builder.WriteString("steps:\n")
	for _, step := range report.Steps {
		builder.WriteString(fmt.Sprintf("- %s\n", step))
	}
	return builder.String()
}
