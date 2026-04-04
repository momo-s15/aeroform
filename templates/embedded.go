// Package templatedata embeds Terraform templates so the CLI works from any working directory.
package templatedata

import "embed"

//go:embed simple pro
var Files embed.FS
