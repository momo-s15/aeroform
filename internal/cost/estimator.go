package cost

import (
	"fmt"
	"sort"
)

type Estimate struct {
	Template     string
	Monthly      float64
	Lines        []string
	OverBudget   bool
	Budget       float64
	HasDomainFee bool
}

func EstimateForSimpleTemplate(template string, hasCustomDomain bool, budget float64) Estimate {
	info, ok := simpleTemplateCosts[template]
	if !ok {
		info = TemplateCost{Monthly: 0.00, Note: "unknown template, cost not estimated"}
	}

	total := info.Monthly
	lines := []string{fmt.Sprintf("%s: $%.2f (%s)", template, info.Monthly, info.Note)}
	if hasCustomDomain {
		total += 0.50
		lines = append(lines, "Route 53 hosted zone: $0.50 (custom domain)")
	}

	if budget <= 0 {
		budget = 20.00
	}

	return Estimate{
		Template:     template,
		Monthly:      total,
		Lines:        lines,
		OverBudget:   total > budget,
		Budget:       budget,
		HasDomainFee: hasCustomDomain,
	}
}

func SupportedSimpleTemplates() []string {
	templates := make([]string, 0, len(simpleTemplateCosts))
	for name := range simpleTemplateCosts {
		templates = append(templates, name)
	}
	sort.Strings(templates)
	return templates
}
