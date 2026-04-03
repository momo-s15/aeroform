package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/providers"
	"github.com/momo-s15/aeroform/internal/security"
	"github.com/momo-s15/aeroform/internal/terraform"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan [prompt]",
	Short: "Pro Mode dry-run: select templates, scan, and show terraform plan without applying",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		printBanner(out)

		cfg, err := config.LoadFromFile("config.yaml")
		if err != nil {
			return fmt.Errorf("pro mode requires config.yaml: %w", err)
		}

		plan, err := engine.BuildProPlanWithClient(cfg, args[0], proLLMClientFromConfig(cfg))
		if err != nil {
			return err
		}

		fmt.Fprintln(out, "Aeroform Pro Mode plan (dry-run)")
		for _, line := range plan.Summary {
			fmt.Fprintln(out, "  "+line)
		}
		fmt.Fprintln(out, "")

		provider := providers.ForCloud(cfg.Cloud)
		vars := provider.GenerateVars(cfg, nil)

		workDir := filepath.Join(".aeroform", "pro", engine.Slugify(args[0]))
		sources := proTemplateSources(provider, plan.Templates)

		estimate := provider.EstimateCost(plan.Templates)
		fmt.Fprintln(out, "Estimated monthly cost:")
		for _, line := range estimate.Lines() {
			fmt.Fprintln(out, line)
		}
		fmt.Fprintln(out, "")

		fmt.Fprintf(out, "-> Rendering %d template(s) into %s\n", len(sources), workDir)
		if err := terraform.RenderProProject(sources, vars, workDir); err != nil {
			return fmt.Errorf("render project: %w", err)
		}

		fmt.Fprintln(out, "-> Running security scan...")
		scanReport := security.RunProSecurityScan(workDir)
		for _, line := range scanReport.SummaryLines() {
			fmt.Fprintln(out, "  "+line)
		}
		if scanReport.Blocking {
			return fmt.Errorf("security gate blocked plan — fix findings above and re-run")
		}

		if err := terraform.EnsureBinary(); err != nil {
			return err
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

		fmt.Fprintln(out, "")
		fmt.Fprintln(out, "Plan-only mode: no infrastructure changes were made.")
		fmt.Fprintln(out, "Run 'aeroform generate' with the same prompt to apply.")

		return nil
	},
}
