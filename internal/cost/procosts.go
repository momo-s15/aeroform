package cost

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type ProTemplateCost struct {
	Monthly decimal.Decimal
	Note    string
}

var freeTierThreshold = decimal.NewFromFloat(1.00)

var awsProCosts = map[string]ProTemplateCost{
	"vpc":            {Monthly: decimal.NewFromFloat(32.40), Note: "NAT Gateway ($32.40/mo baseline) + data transfer"},
	"eks":            {Monthly: decimal.NewFromFloat(73.00), Note: "EKS control plane ($73/mo) + node instances extra"},
	"rds-private":    {Monthly: decimal.NewFromFloat(15.00), Note: "db.t3.micro with 20GB gp3"},
	"s3-private":     {Monthly: decimal.NewFromFloat(0.50), Note: "storage + requests at low volume"},
	"alb":            {Monthly: decimal.NewFromFloat(22.00), Note: "ALB fixed ($16.20) + LCU cost"},
	"lambda-api":     {Monthly: decimal.NewFromInt(0), Note: "free tier covers 1M requests/mo"},
	"ecs-fargate":    {Monthly: decimal.NewFromFloat(36.00), Note: "0.25 vCPU / 0.5GB per task"},
	"cloudfront-api": {Monthly: decimal.NewFromFloat(1.00), Note: "first 1TB free, WAF $5/mo if enabled"},
}

var azureProCosts = map[string]ProTemplateCost{
	"vnet":        {Monthly: decimal.NewFromInt(0), Note: "VNet is free; peering and NAT extra"},
	"aks":         {Monthly: decimal.NewFromFloat(73.00), Note: "AKS free tier control plane + node VM cost"},
	"cosmos-db":   {Monthly: decimal.NewFromFloat(25.00), Note: "serverless 400 RU/s baseline"},
	"app-service": {Monthly: decimal.NewFromFloat(13.14), Note: "B1 plan Linux"},
	"storage":     {Monthly: decimal.NewFromFloat(1.00), Note: "LRS hot storage at low volume"},
	"key-vault":   {Monthly: decimal.NewFromFloat(0.03), Note: "$0.03 per 10K operations"},
}

var gcpProCosts = map[string]ProTemplateCost{
	"vpc":      {Monthly: decimal.NewFromFloat(32.40), Note: "Cloud NAT baseline + data processing"},
	"gke":      {Monthly: decimal.NewFromFloat(73.00), Note: "GKE Autopilot/Standard + node cost"},
	"cloudsql": {Monthly: decimal.NewFromFloat(7.67), Note: "db-f1-micro shared-core + 10GB SSD"},
	"gcs":      {Monthly: decimal.NewFromFloat(0.50), Note: "standard storage at low volume"},
}

func proTableForCloud(cloud string) map[string]ProTemplateCost {
	switch cloud {
	case "azure":
		return azureProCosts
	case "gcp":
		return gcpProCosts
	default:
		return awsProCosts
	}
}

type ProEstimate struct {
	Cloud    string
	Items    []ProLineItem
	Monthly  decimal.Decimal
	FreeTier bool
}

type ProLineItem struct {
	Template string
	Monthly  decimal.Decimal
	Note     string
}

func EstimateForProTemplates(cloud string, templates []string) ProEstimate {
	table := proTableForCloud(cloud)
	est := ProEstimate{Cloud: cloud, Monthly: decimal.NewFromInt(0)}
	for _, t := range templates {
		info, ok := table[t]
		if !ok {
			info = ProTemplateCost{Monthly: decimal.NewFromInt(0), Note: "cost not yet catalogued"}
		}
		est.Items = append(est.Items, ProLineItem{
			Template: t,
			Monthly:  info.Monthly,
			Note:     info.Note,
		})
		est.Monthly = est.Monthly.Add(info.Monthly)
	}
	est.FreeTier = est.Monthly.LessThan(freeTierThreshold)
	return est
}

func (e ProEstimate) Lines() []string {
	lines := make([]string, 0, len(e.Items)+2)
	for _, item := range e.Items {
		lines = append(lines, fmt.Sprintf("  %-18s $%7s/mo  (%s)", item.Template, item.Monthly.StringFixed(2), item.Note))
	}
	lines = append(lines, fmt.Sprintf("  %-18s $%7s/mo", "TOTAL (estimated)", e.Monthly.StringFixed(2)))
	if e.FreeTier {
		lines = append(lines, "  All selected services fall within free tier limits.")
	}
	return lines
}
