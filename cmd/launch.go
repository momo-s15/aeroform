package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/llm"
	"github.com/momo-s15/aeroform/internal/logger"
	"github.com/momo-s15/aeroform/internal/security"
	"github.com/momo-s15/aeroform/internal/simplestate"
	"github.com/momo-s15/aeroform/internal/terraform"
	"github.com/momo-s15/aeroform/internal/ui"
	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch a guided Simple Mode deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.L()
		mode := engine.DetectMode("config.yaml", engine.Flags{})
		if mode == engine.ProMode {
			return errors.New("launch is for Simple Mode; use 'aeroform generate' for Pro Mode")
		}

		out := cmd.OutOrStdout()
		printBanner(out)

		plan, err := gatherSimpleLaunchPlan(out)
		if err != nil {
			return err
		}

		fmt.Fprintln(out, "→ Checking plan against Simple Mode security rules…")
		report := security.EvaluateSimplePlan(plan)
		security.PrintSimpleReport(out, report)
		if report.BlockingCount > 0 {
			return errors.New("simple mode security gate blocked this launch; fix the issues above and try again")
		}
		plan = report.CorrectedPlan

		printSimpleLaunchSummary(out, plan)

		proceed, err := uiConfirm("Continue with deployment")
		if err != nil {
			return err
		}
		if !proceed {
			fmt.Fprintln(out, "Launch cancelled.")
			return nil
		}

		if err := terraform.EnsureBinary(); err != nil {
			return err
		}

		workDir := filepath.Join(".aeroform", "projects", plan.ProjectName)
		log.Debugw("rendering template", "template", plan.Template, "workDir", workDir)
		fmt.Fprintf(out, "\n-> Rendering template %s into %s\n", plan.Template, workDir)
		if err := terraform.RenderTemplate(plan.TemplateDir, nil, workDir); err != nil {
			return fmt.Errorf("render template: %w", err)
		}
		if err := checkAzureSimpleTemplateFresh(workDir, plan); err != nil {
			return err
		}

		vars := simpleTfvars(plan)
		if err := terraform.GenerateTfvars(vars, workDir); err != nil {
			return fmt.Errorf("generate tfvars: %w", err)
		}
		if plan.Provider == "azure" {
			fmt.Fprintf(out, "-> Using Azure location %q for this deploy (if this is wrong, set AEROFORM_AZURE_LOCATION and run again)\n", vars["location"])
		}

		corrections, err := security.AutoCorrect(workDir)
		if err != nil {
			return fmt.Errorf("security autocorrect: %w", err)
		}
		if len(corrections) > 0 {
			fmt.Fprintln(out, "-> Security auto-corrections applied:")
			for _, c := range corrections {
				fmt.Fprintf(out, "   [%s] %s (%s)\n", c.Rule, c.Description, c.File)
			}
		}

		fmt.Fprintln(out, "-> Running terraform init...")
		if err := terraform.Init(workDir); err != nil {
			return err
		}

		fmt.Fprintln(out, "-> Running terraform plan...")
		planResult, err := terraform.Plan(workDir)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "   Plan: %d to add, %d to change, %d to destroy.\n",
			planResult.AddCount, planResult.ChangeCount, planResult.DestroyCount)

		applyOk, err := uiConfirm("Apply this plan")
		if err != nil {
			return err
		}
		if !applyOk {
			fmt.Fprintln(out, "Apply cancelled. Your plan is saved — run 'aeroform launch' again to resume.")
			return nil
		}

		fmt.Fprintln(out, "-> Deploying...")
		log.Debugw("running terraform apply", "workDir", workDir)
		if err := terraform.Apply(workDir, true); err != nil {
			return err
		}
		log.Debug("terraform apply completed successfully")

		if err := simplestate.AddProject(simplestate.Project{
			Name:            plan.ProjectName,
			Provider:        plan.Provider,
			Template:        plan.Template,
			Prompt:          plan.Prompt,
			CustomDomain:    plan.CustomDomain,
			MonthlyEstimate: plan.Cost.Monthly.InexactFloat64(),
			WorkDir:         workDir,
			CreatedAt:       time.Now().UTC(),
		}); err != nil {
			return err
		}

		fmt.Fprintln(out, "")
		ui.Successln(out, "✓ Done! Your infrastructure is live.")
		outputs, err := terraform.Output(workDir)
		if err == nil && outputs != "" {
			fmt.Fprintln(out, "")
			ui.Boldln(out, "Outputs:")
			fmt.Fprintln(out, outputs)
		}

		return nil
	},
}

func gatherSimpleLaunchPlan(out io.Writer) (engine.SimpleLaunchPlan, error) {
	request, err := uiPrompt("What do you want to launch", "a personal website", validateNotEmpty)
	if err != nil {
		return engine.SimpleLaunchPlan{}, err
	}

	providerChoices := []string{"aws", "azure", "gcp"}
	selected, err := uiSelect("Cloud provider", providerChoices)
	if err != nil {
		return engine.SimpleLaunchPlan{}, err
	}
	provider := strings.Fields(selected)[0]

	defaultSlug := engine.Slugify(request)
	if validateSlug(defaultSlug) != nil {
		defaultSlug = "my-project"
	}
	projectName, err := uiPrompt("Project name", defaultSlug, validateSlug)
	if err != nil {
		return engine.SimpleLaunchPlan{}, err
	}

	customDomain, err := uiPromptOptional("Custom domain (leave blank to skip)")
	if err != nil {
		return engine.SimpleLaunchPlan{}, err
	}

	var azureLocation string
	if provider == "azure" {
		azDef := firstNonEmpty(
			os.Getenv("AEROFORM_AZURE_LOCATION"),
			os.Getenv("AZURE_LOCATION"),
			os.Getenv("AZURE_DEFAULT_REGION"),
			"canadacentral",
		)
		azureLocation, err = uiPrompt("Azure region (try canadacentral on Azure for Education if eastus is blocked)", azDef, validateAzureRegion)
		if err != nil {
			return engine.SimpleLaunchPlan{}, err
		}
	}

	var gcpProjectID string
	if provider == "gcp" {
		gcpDefault := firstNonEmpty(os.Getenv("GOOGLE_PROJECT"), os.Getenv("GCP_PROJECT"), os.Getenv("CLOUDSDK_CORE_PROJECT"))
		gcpProjectID, err = uiPrompt("GCP project ID", gcpDefault, validateGCPProjectID)
		if err != nil {
			return engine.SimpleLaunchPlan{}, err
		}
	}

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "→ Matching your request to a template (local AI may take a few seconds)…")

	plan, err := engine.BuildSimpleLaunchPlanWithClient(engine.SimpleLaunchInput{
		Prompt:        request,
		Provider:      provider,
		ProjectName:   projectName,
		CustomDomain:  customDomain,
		GCPProjectID:  gcpProjectID,
		AzureLocation: azureLocation,
	}, llmClientForSimpleMode())
	if err != nil {
		return engine.SimpleLaunchPlan{}, err
	}

	return plan, nil
}

func llmClientForSimpleMode() llm.Client {
	client := engine.DefaultSimpleLLMClient()
	if client != nil && client.IsAvailable() {
		return client
	}
	return nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

// simpleTfvars maps provider-specific Terraform variable names (AWS uses region; Azure uses location; GCP needs project_id).
func simpleTfvars(plan engine.SimpleLaunchPlan) map[string]string {
	vars := map[string]string{"project_name": plan.ProjectName}
	if plan.CustomDomain != "" {
		vars["custom_domain"] = plan.CustomDomain
	}
	switch plan.Provider {
	case "aws":
		vars["region"] = firstNonEmpty(os.Getenv("AWS_REGION"), os.Getenv("AWS_DEFAULT_REGION"), "us-east-1")
	case "azure":
		vars["location"] = firstNonEmpty(plan.AzureLocation, os.Getenv("AEROFORM_AZURE_LOCATION"), os.Getenv("AZURE_LOCATION"), "canadacentral")
	case "gcp":
		vars["region"] = firstNonEmpty(os.Getenv("GOOGLE_REGION"), os.Getenv("GCP_REGION"), "us-central1")
		vars["project_id"] = plan.GCPProjectID
	}
	return vars
}

func printSimpleLaunchSummary(out io.Writer, plan engine.SimpleLaunchPlan) {
	fmt.Fprintln(out, "")
	ui.Boldln(out, "Aeroform Simple Mode launch")
	for _, line := range plan.Summary {
		fmt.Fprintln(out, "  "+line)
	}
	fmt.Fprintln(out, "")
	ui.Boldln(out, "Estimated monthly cost")
	for _, line := range plan.Cost.Lines {
		fmt.Fprintln(out, "  "+line)
	}
	if plan.Cost.Monthly.IsZero() {
		ui.Success(out, "  total: $%s/month (free tier)\n", plan.Cost.Monthly.StringFixed(2))
	} else {
		ui.Warn(out, "  total: $%s/month\n", plan.Cost.Monthly.StringFixed(2))
	}
	if plan.Cost.OverBudget {
		ui.Warn(out, "  ⚠ warning: this is over the $%s default budget\n", plan.Cost.Budget.StringFixed(2))
	}
	fmt.Fprintln(out, "")
}

// Marker in templates/simple/azure/static-site/main.tf — detects stale CLI or hand-edited copies.
const azureStaticSiteSchemaMarker = "aeroform-schema: azure-static-site/3"

func checkAzureSimpleTemplateFresh(workDir string, plan engine.SimpleLaunchPlan) error {
	if plan.Provider != "azure" || plan.Template != "static-site" {
		return nil
	}
	b, err := os.ReadFile(filepath.Join(workDir, "main.tf"))
	if err != nil {
		return nil
	}
	s := string(b)
	// Stale copies used regexreplace(), which many Terraform builds reject; catch before plan.
	if strings.Contains(s, "regexreplace(") {
		return fmt.Errorf(
			"rendered Azure static-site main.tf in %q still calls regexreplace() (stale aeroform build, or AEROFORM_DISK_TEMPLATES=1 with old files under .\\templates\\).\n"+
				"(this binary: aeroform %s)\n\n"+
				"Fix: go install github.com/momo-s15/aeroform@main, run: Remove-Item Env:AEROFORM_DISK_TEMPLATES -ErrorAction SilentlyContinue,\n"+
				"delete folder %q, then run launch again",
			workDir, Version, workDir,
		)
	}
	if strings.Contains(s, azureStaticSiteSchemaMarker) {
		return nil
	}
	return fmt.Errorf(
		"Azure static-site files in %q are missing %q (stale aeroform binary, old templates/ checkout, or a project folder from an older run).\n"+
			"That often defaults the region to eastus and can destroy/recreate your resource group.\n\n"+
			"Fix: go install github.com/momo-s15/aeroform@main (or git pull if you build from source), remove that project folder under .aeroform/projects, and run launch again",
		workDir, azureStaticSiteSchemaMarker,
	)
}
