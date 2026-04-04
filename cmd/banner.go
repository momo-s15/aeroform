package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func bannerVersionLabel() string {
	v := Version
	if v == "dev" {
		return "dev"
	}
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}

// printBanner prints the Aeroform ASCII wordmark, feature table, and taglines.
// Honors NO_COLOR; outputs to w (use cmd.OutOrStdout() from commands).
func printBanner(w io.Writer) {
	noColor := func() bool {
		_, ok := os.LookupEnv("NO_COLOR")
		return ok
	}()

	c := func(code, text string) string {
		if noColor {
			return text
		}
		return code + text + "\033[0m"
	}

	// Each row gets a progressively darker cyan → deep blue shade
	rows := []struct {
		color string
		text  string
	}{
		{"\033[38;2;41;171;226m", ` █████╗ ███████╗██████╗  ██████╗ ███████╗ ██████╗ ██████╗ ███╗   ███╗`},
		{"\033[38;2;35;155;208m", `██╔══██╗██╔════╝██╔══██╗██╔═══██╗██╔════╝██╔═══██╗██╔══██╗████╗ ████║`},
		{"\033[38;2;26;135;185m", `███████║█████╗  ██████╔╝██║   ██║█████╗  ██║   ██║██████╔╝██╔████╔██║`},
		{"\033[38;2;18;112;160m", `██╔══██║██╔══╝  ██╔══██╗██║   ██║██╔══╝  ██║   ██║██╔══██╗██║╚██╔╝██║`},
		{"\033[38;2;12;88;130m", `██║  ██║███████╗██║  ██║╚██████╔╝██║     ╚██████╔╝██║  ██║██║ ╚═╝ ██║`},
		{"\033[38;2;7;60;90m", `╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝      ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝`},
	}

	for _, row := range rows {
		fmt.Fprintln(w, c(row.color, row.text))
	}

	sep := "\033[38;2;41;171;226m"
	wht := "\033[97m"
	grn := "\033[38;2;63;185;80m"
	gray := "\033[90m"
	reset := "\033[0m"

	if noColor {
		sep, wht, grn, gray, reset = "", "", "", "", ""
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s ┌─────────────────────┬──────────────────────┬─────────────────────┐%s\n", sep, reset)
	fmt.Fprintf(w, "%s│%s  %-21s%s│%s  %-21s%s│%s  %-21s%s│\n",
		sep, wht, "clouds", reset, wht, "AI backend", reset, wht, "security", reset)
	fmt.Fprintf(w, "%s│%s  %-21s%s│%s  %-21s%s│%s  %-21s%s│\n",
		sep, grn, "AWS · Azure · GCP", reset, grn, "Ollama (free, local)", reset, grn, "Checkov + tfsec", reset)
	fmt.Fprintf(w, "%s └─────────────────────┴──────────────────────┴─────────────────────┘%s\n", sep, reset)

	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s  cloud infrastructure from plain english%s\n", wht, reset)
	fmt.Fprintf(w, "%s  %s  ·  free AI  ·  github.com/momo-s15/aeroform%s\n\n", gray, bannerVersionLabel(), reset)
}
