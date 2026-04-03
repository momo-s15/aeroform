package aws

import (
	"github.com/momo-s15/aeroform/bootstrap"
	"github.com/momo-s15/aeroform/internal/config"
)

func Run(cfg config.Config) bootstrap.Report {
	return bootstrap.Run(cfg)
}
