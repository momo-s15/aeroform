package prompt

import (
	"fmt"
	"strings"
)

func BuildProPrompt(userPrompt, cloud, region, environment string, allowedTemplates []string) string {
	if environment == "" {
		environment = "production"
	}

	return fmt.Sprintf(`You are an IaC assistant for professional cloud engineers.
Choose production-grade resources appropriate to the request.

Available pro templates: [%s]

User config: cloud=%s region=%s env=%s
User: %q

Respond ONLY with JSON:
{"templates":["vpc","eks","rds-private"],"params":{"environment":"%s"}}`,
		strings.Join(allowedTemplates, ", "),
		cloud,
		region,
		environment,
		userPrompt,
		environment,
	)
}
