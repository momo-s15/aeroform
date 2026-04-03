package engine

import "testing"

func TestSelectSimpleTemplate(t *testing.T) {
	tests := []struct {
		prompt string
		want   string
	}{
		{prompt: "a contact form for my site", want: "contact-form"},
		{prompt: "a backend api", want: "lambda-api"},
		{prompt: "a database for my project", want: "tiny-db"},
		{prompt: "a personal website", want: "static-site"},
	}

	for _, testCase := range tests {
		t.Run(testCase.prompt, func(t *testing.T) {
			if got := SelectSimpleTemplate(testCase.prompt); got != testCase.want {
				t.Fatalf("SelectSimpleTemplate(%q) = %q, want %q", testCase.prompt, got, testCase.want)
			}
		})
	}
}

func TestSelectSimpleTemplateAzure(t *testing.T) {
	tests := []struct {
		prompt string
		want   string
	}{
		{prompt: "a backend api", want: "function-api"},
		{prompt: "a serverless function", want: "function-api"},
		{prompt: "a personal website", want: "static-site"},
	}

	for _, testCase := range tests {
		t.Run(testCase.prompt, func(t *testing.T) {
			if got := SelectSimpleTemplateForProvider(testCase.prompt, "azure"); got != testCase.want {
				t.Fatalf("SelectSimpleTemplateForProvider(%q, azure) = %q, want %q", testCase.prompt, got, testCase.want)
			}
		})
	}
}

func TestSelectSimpleTemplateGCP(t *testing.T) {
	tests := []struct {
		prompt string
		want   string
	}{
		{prompt: "a backend api", want: "cloud-run-api"},
		{prompt: "a container service", want: "cloud-run-api"},
		{prompt: "a personal website", want: "static-site"},
	}

	for _, testCase := range tests {
		t.Run(testCase.prompt, func(t *testing.T) {
			if got := SelectSimpleTemplateForProvider(testCase.prompt, "gcp"); got != testCase.want {
				t.Fatalf("SelectSimpleTemplateForProvider(%q, gcp) = %q, want %q", testCase.prompt, got, testCase.want)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	if got := Slugify("My Personal Website!"); got != "my-personal-website" {
		t.Fatalf("Slugify() = %q", got)
	}
}

func TestBuildSimpleLaunchPlan(t *testing.T) {
	plan, err := BuildSimpleLaunchPlan(SimpleLaunchInput{
		Prompt:       "a personal website",
		Provider:     "aws",
		ProjectName:  "",
		CustomDomain: "example.com",
	})
	if err != nil {
		t.Fatalf("BuildSimpleLaunchPlan() error = %v", err)
	}
	if plan.Template != "static-site" {
		t.Fatalf("unexpected template %q", plan.Template)
	}
	if plan.ProjectName != "a-personal-website" {
		t.Fatalf("unexpected project name %q", plan.ProjectName)
	}
}

func TestBuildSimpleLaunchPlanAzure(t *testing.T) {
	plan, err := BuildSimpleLaunchPlan(SimpleLaunchInput{
		Prompt:   "a personal website",
		Provider: "azure",
	})
	if err != nil {
		t.Fatalf("BuildSimpleLaunchPlan(azure) error = %v", err)
	}
	if plan.Provider != "azure" {
		t.Fatalf("expected provider=azure, got %q", plan.Provider)
	}
	if plan.Template != "static-site" {
		t.Fatalf("expected static-site for azure website, got %q", plan.Template)
	}
}

func TestBuildSimpleLaunchPlanGCP(t *testing.T) {
	plan, err := BuildSimpleLaunchPlan(SimpleLaunchInput{
		Prompt:   "a backend api",
		Provider: "gcp",
	})
	if err != nil {
		t.Fatalf("BuildSimpleLaunchPlan(gcp) error = %v", err)
	}
	if plan.Provider != "gcp" {
		t.Fatalf("expected provider=gcp, got %q", plan.Provider)
	}
	if plan.Template != "cloud-run-api" {
		t.Fatalf("expected cloud-run-api for gcp api, got %q", plan.Template)
	}
}

func TestBuildSimpleLaunchPlanUnsupported(t *testing.T) {
	_, err := BuildSimpleLaunchPlan(SimpleLaunchInput{
		Prompt:   "a website",
		Provider: "digitalocean",
	})
	if err == nil {
		t.Fatal("expected error for unsupported provider")
	}
}
