package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/momo-s15/aeroform/internal/engine"
	"github.com/momo-s15/aeroform/internal/llm"
	"github.com/momo-s15/aeroform/internal/security"
	"github.com/momo-s15/aeroform/internal/simplestate"
	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch a guided Simple Mode deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		mode := engine.DetectMode("config.yaml", engine.Flags{})
		if mode == engine.ProMode {
			return errors.New("pro mode is not implemented yet; remove config.yaml or wait for the next phase")
		}

		plan, err := gatherSimpleLaunchPlan(cmd)
		if err != nil {
			return err
		}

		report := security.EvaluateSimplePlan(plan)
		security.PrintSimpleReport(cmd.OutOrStdout(), report)
		if report.BlockingCount > 0 {
			return errors.New("simple mode security gate blocked this launch; fix the issues above and try again")
		}

		plan = report.CorrectedPlan
		if err := simplestate.AddProject(simplestate.Project{
			Name:            plan.ProjectName,
			Provider:        plan.Provider,
			Template:        plan.Template,
			Prompt:          plan.Prompt,
			CustomDomain:    plan.CustomDomain,
			MonthlyEstimate: plan.Cost.Monthly,
			CreatedAt:       time.Now().UTC(),
		}); err != nil {
			return err
		}

		return printSimpleLaunchPlan(cmd.OutOrStdout(), plan)
	},
}

func gatherSimpleLaunchPlan(cmd *cobra.Command) (engine.SimpleLaunchPlan, error) {
	reader := bufio.NewReader(cmd.InOrStdin())
	out := cmd.OutOrStdout()

	request, err := askPrompt(reader, out, "What do you want to launch?", "a personal website")
	if err != nil {
		return engine.SimpleLaunchPlan{}, err
	}

	provider, err := askPrompt(reader, out, "Which cloud provider do you want to use? [aws available now]", "aws")
	if err != nil {
		return engine.SimpleLaunchPlan{}, err
	}

	projectName, err := askPrompt(reader, out, "What should we call this project?", engine.Slugify(request))
	if err != nil {
		return engine.SimpleLaunchPlan{}, err
	}

	customDomain, err := askOptionalPrompt(reader, out, "Custom domain (optional)")
	if err != nil {
		return engine.SimpleLaunchPlan{}, err
	}

	plan, err := engine.BuildSimpleLaunchPlanWithClient(engine.SimpleLaunchInput{
		Prompt:       request,
		Provider:     provider,
		ProjectName:  projectName,
		CustomDomain: customDomain,
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

func askPrompt(reader *bufio.Reader, out io.Writer, question string, defaultValue string) (string, error) {
	fmt.Fprintf(out, "%s\n", question)
	if defaultValue != "" {
		fmt.Fprintf(out, "[%s] ", defaultValue)
	} else {
		fmt.Fprint(out, "> ")
	}

	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultValue, nil
	}
	return line, nil
}

func askOptionalPrompt(reader *bufio.Reader, out io.Writer, question string) (string, error) {
	fmt.Fprintf(out, "%s\n", question)
	fmt.Fprint(out, "> ")
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func printSimpleLaunchPlan(out io.Writer, plan engine.SimpleLaunchPlan) error {
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Aeroform Simple Mode launch")
	for _, line := range plan.Summary {
		fmt.Fprintln(out, line)
	}
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Estimated monthly cost")
	for _, line := range plan.Cost.Lines {
		fmt.Fprintln(out, line)
	}
	fmt.Fprintf(out, "total: $%.2f/month\n", plan.Cost.Monthly)
	if plan.Cost.OverBudget {
		fmt.Fprintf(out, "warning: this is over the $%.2f default budget\n", plan.Cost.Budget)
	}
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Next step")
	fmt.Fprintf(out, "- template dir: %s\n", plan.TemplateDir)
	fmt.Fprintf(out, "- project name: %s\n", plan.ProjectName)
	if plan.CustomDomain != "" {
		fmt.Fprintf(out, "- custom domain: %s\n", plan.CustomDomain)
	}
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Deployment is not wired yet; Phase 2 is now selecting templates and producing a costed launch plan.")
	return nil
}
