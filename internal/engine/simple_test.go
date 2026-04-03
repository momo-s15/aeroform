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
