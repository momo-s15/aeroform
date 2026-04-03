package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/momo-s15/aeroform/internal/config"
	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/llm"
	"github.com/momo-s15/aeroform/internal/logger"
	"github.com/momo-s15/aeroform/internal/providers"
	"github.com/momo-s15/aeroform/internal/security"
	"github.com/momo-s15/aeroform/internal/terraform"
	"github.com/momo-s15/aeroform/internal/ui"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [prompt]",
	Short: "Pro Mode pipeline: select templates, scan, plan, and deploy",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		printBanner(cmd.OutOrStdout())
		log := logger.L()
		log.Debugw("loading config", "path", "config.yaml")
		cfg, err := config.LoadFromFile("config.yaml")
		if err != nil {
			return fmt.Errorf("pro mode requires config.yaml: %w", err)
		}

		log.Debugw("building pro plan", "prompt", args[0], "cloud", cfg.Cloud)
		plan, err := engine.BuildProPlanWithClient(cfg, args[0], proLLMClientFromConfig(cfg))
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		ui.Boldln(out, "Aeroform Pro Mode generate")
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
			ui.Errorln(out, "✗ Security gate blocked generate — fix findings above and re-run")
			return fmt.Errorf("security gate blocked generate — fix findings above and re-run")
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

		proceed, err := uiConfirm("Apply this plan")
		if err != nil {
			return err
		}
		if !proceed {
			fmt.Fprintln(out, "Apply cancelled. Run 'aeroform generate' again to resume.")
			return nil
		}

		fmt.Fprintln(out, "-> Deploying...")
		log.Debugw("running terraform apply", "workDir", workDir)
		if err := terraform.Apply(workDir, true); err != nil {
			return err
		}

		log.Debug("terraform apply completed successfully")
		fmt.Fprintln(out, "")
		ui.Successln(out, "✓ Done! Pro Mode infrastructure is live.")
		outputs, err := terraform.Output(workDir)
		if err == nil && outputs != "" {
			fmt.Fprintln(out, "")
			ui.Boldln(out, "Outputs:")
			fmt.Fprintln(out, outputs)
		}

		fmt.Fprintln(out, "")
		if err := provider.PostDeploy(cfg, out); err != nil {
			ui.Warn(out, "⚠ Warning: post-deploy steps failed: %v\n", err)
		}

		return nil
	},
}

func proLLMClientFromConfig(cfg config.Config) llm.Client {
	client := engine.ProLLMClientFromConfig(cfg)
	if client != nil && client.IsAvailable() {
		return client
	}
	return nil
}

func proTemplateSources(provider providers.CloudProvider, templates []string) []terraform.ProTemplateSource {
	sources := make([]terraform.ProTemplateSource, 0, len(templates))
	for _, t := range templates {
		sources = append(sources, terraform.ProTemplateSource{
			Name: t,
			Dir:  provider.GetTemplateDir(engine.ProMode, t),
		})
	}
	return sources
}
