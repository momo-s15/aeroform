package prompt

import (
	"fmt"
	"strings"
)

func BuildSimplePrompt(userPrompt, cloud, region string, allowedTemplates []string) string {
	return fmt.Sprintf(`You are an infrastructure assistant for student developers.
PRIORITY: Choose free-tier resources whenever possible.
NEVER choose: EC2 > t3.micro, RDS multi-AZ, EKS, NAT Gateway.

Available simple templates: [%s]

User: %q
Cloud: %s
Region: %s

Respond ONLY with JSON:
{"templates":["static-site"],"params":{"project_name":"..."}}`, strings.Join(allowedTemplates, ", "), userPrompt, cloud, region)
}
