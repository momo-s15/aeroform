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

// printBanner prints the Aeroform ASCII wordmark and taglines (cyan / white / gray).
// Honors NO_COLOR; outputs to w (use cmd.OutOrStdout() from commands).
func printBanner(w io.Writer) {
	cyan := "\033[36m"
	white := "\033[97m"
	gray := "\033[90m"
	reset := "\033[0m"
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		cyan, white, gray, reset = "", "", "", ""
	}

	bannerArt := cyan + ` █████╗ ███████╗██████╗  ██████╗ ███████╗ ██████╗ ██████╗ ███╗   ███╗
██╔══██╗██╔════╝██╔══██╗██╔═══██╗██╔════╝██╔═══██╗██╔══██╗████╗ ████║
███████║█████╗  ██████╔╝██║   ██║█████╗  ██║   ██║██████╔╝██╔████╔██║
██╔══██║██╔══╝  ██╔══██╗██║   ██║██╔══╝  ██║   ██║██╔══██╗██║╚██╔╝██║
██║  ██║███████╗██║  ██║╚██████╔╝██║     ╚██████╔╝██║  ██║██║ ╚═╝ ██║
╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝      ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝
` + reset
	fmt.Fprint(w, bannerArt)
	fmt.Fprint(w, white+"  cloud infrastructure from plain english  ·  AWS  ·  Azure  ·  GCP\n"+reset)
	fmt.Fprintf(w, "%s  %s  ·  free AI  ·  github.com/momo-s15/aeroform%s\n\n", gray, bannerVersionLabel(), reset)
}
