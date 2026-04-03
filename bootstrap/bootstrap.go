package bootstrap

import (
	"fmt"
	"strings"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/providers"
)

type Report struct {
	Cloud    string
	Repo     string
	Sections []providers.BootstrapSection
}

type Params struct {
	Repo string
}

func Run(cfg config.Config, params Params) Report {
	provider := providers.ForCloud(cfg.Cloud)
	sections := provider.BootstrapSections(cfg, params.Repo)

	return Report{
		Cloud:    provider.Name(),
		Repo:     params.Repo,
		Sections: sections,
	}
}

func Render(report Report) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Aeroform bootstrap — %s\n", report.Cloud))
	if report.Repo != "" {
		b.WriteString(fmt.Sprintf("GitHub repo: %s\n", report.Repo))
	}
	b.WriteString("\n")

	for i, section := range report.Sections {
		b.WriteString(fmt.Sprintf("=== %d. %s ===\n\n", i+1, section.Title))
		for _, cmd := range section.Commands {
			b.WriteString(cmd)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}
