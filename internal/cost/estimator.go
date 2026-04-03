package cost

import (
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
)

var domainFee = decimal.NewFromFloat(0.50)
var defaultBudget = decimal.NewFromFloat(20.00)

type Estimate struct {
	Template     string
	Monthly      decimal.Decimal
	Lines        []string
	OverBudget   bool
	Budget       decimal.Decimal
	HasDomainFee bool
}

func EstimateForSimpleTemplate(provider, template string, hasCustomDomain bool, budget float64) Estimate {
	info, ok := lookupSimpleCost(provider, template)
	if !ok {
		info = TemplateCost{Monthly: decimal.NewFromInt(0), Note: "unknown template, cost not estimated"}
	}

	total := info.Monthly
	lines := []string{fmt.Sprintf("%s: $%s (%s)", template, info.Monthly.StringFixed(2), info.Note)}
	if hasCustomDomain {
		total = total.Add(domainFee)
		lines = append(lines, "DNS hosted zone: $0.50 (custom domain)")
	}

	budgetDec := decimal.NewFromFloat(budget)
	if budgetDec.LessThanOrEqual(decimal.Zero) {
		budgetDec = defaultBudget
	}

	return Estimate{
		Template:     template,
		Monthly:      total,
		Lines:        lines,
		OverBudget:   total.GreaterThan(budgetDec),
		Budget:       budgetDec,
		HasDomainFee: hasCustomDomain,
	}
}

func lookupSimpleCost(provider, template string) (TemplateCost, bool) {
	if provider == "" {
		provider = "aws"
	}
	table, ok := simpleTemplateCosts[provider]
	if !ok {
		return TemplateCost{}, false
	}
	info, ok := table[template]
	return info, ok
}

func SupportedSimpleTemplates() []string {
	return SupportedSimpleTemplatesForProvider("aws")
}

func SupportedSimpleTemplatesForProvider(provider string) []string {
	if provider == "" {
		provider = "aws"
	}
	table, ok := simpleTemplateCosts[provider]
	if !ok {
		return nil
	}
	templates := make([]string, 0, len(table))
	for name := range table {
		templates = append(templates, name)
	}
	sort.Strings(templates)
	return templates
}
