package cost

import "github.com/shopspring/decimal"

type TemplateCost struct {
	Monthly decimal.Decimal
	Note    string
}

var simpleTemplateCosts = map[string]map[string]TemplateCost{
	"aws": {
		"static-site": {
			Monthly: decimal.NewFromFloat(0.50),
			Note:    "S3 and CloudFront are free tier friendly; Route 53 adds a small hosted zone cost",
		},
		"contact-form": {
			Monthly: decimal.NewFromFloat(0.50),
			Note:    "Serverless contact form with a small Route 53 cost if you add a custom domain",
		},
		"lambda-api": {
			Monthly: decimal.NewFromInt(0),
			Note:    "Lambda, API Gateway, and DynamoDB are typically free tier friendly at small scale",
		},
		"tiny-db": {
			Monthly: decimal.NewFromFloat(14.99),
			Note:    "Private RDS PostgreSQL t3.micro is a low-cost starter database",
		},
		"discord-bot": {
			Monthly: decimal.NewFromFloat(8.00),
			Note:    "EC2 t3.micro with Elastic IP for a persistent bot host",
		},
		"game-server": {
			Monthly: decimal.NewFromFloat(30.00),
			Note:    "EC2 t3.medium with 30GB EBS for game workloads",
		},
		"fullstack-app": {
			Monthly: decimal.NewFromFloat(15.00),
			Note:    "App Runner with private RDS t3.micro and S3 for assets",
		},
		"file-upload": {
			Monthly: decimal.NewFromInt(0),
			Note:    "S3 + Lambda + API Gateway with presigned URLs are free tier friendly",
		},
		"url-shortener": {
			Monthly: decimal.NewFromInt(0),
			Note:    "Lambda + API Gateway + DynamoDB are free tier friendly at small scale",
		},
		"cron-job": {
			Monthly: decimal.NewFromInt(0),
			Note:    "Lambda + EventBridge scheduled rule is 100% free tier",
		},
	},
	"azure": {
		"static-site": {
			Monthly: decimal.NewFromFloat(0.50),
			Note:    "Storage Account static website (HTTPS); classic CDN creation no longer allowed on new Azure subscriptions",
		},
		"function-api": {
			Monthly: decimal.NewFromInt(0),
			Note:    "Azure Functions consumption plan is free tier friendly at small scale",
		},
	},
	"gcp": {
		"static-site": {
			Monthly: decimal.NewFromFloat(0.50),
			Note:    "GCS bucket + Cloud CDN + HTTP load balancer at low traffic",
		},
		"cloud-run-api": {
			Monthly: decimal.NewFromInt(0),
			Note:    "Cloud Run scales to zero; free tier covers 2M requests/mo",
		},
	},
}
