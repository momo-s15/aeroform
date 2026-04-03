package cost

type TemplateCost struct {
	Monthly float64
	Note    string
}

var simpleTemplateCosts = map[string]TemplateCost{
	"static-site": {
		Monthly: 0.50,
		Note:    "S3 and CloudFront are free tier friendly; Route 53 adds a small hosted zone cost",
	},
	"contact-form": {
		Monthly: 0.50,
		Note:    "Serverless contact form with a small Route 53 cost if you add a custom domain",
	},
	"lambda-api": {
		Monthly: 0.00,
		Note:    "Lambda, API Gateway, and DynamoDB are typically free tier friendly at small scale",
	},
	"tiny-db": {
		Monthly: 14.99,
		Note:    "Private RDS PostgreSQL t3.micro is a low-cost starter database",
	},
}
