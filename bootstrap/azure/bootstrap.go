package azure

import (
	"github.com/momo-s15/aeroform/bootstrap"
	"github.com/momo-s15/aeroform/internal/config"
)

func Run(cfg config.Config, params bootstrap.Params) bootstrap.Report {
	return bootstrap.Run(cfg, params)
}
