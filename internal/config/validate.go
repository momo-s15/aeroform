package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterStructValidation(crossFieldValidation, Config{})
}

func crossFieldValidation(sl validator.StructLevel) {
	cfg, ok := sl.Current().Interface().(Config)
	if !ok {
		return
	}
	switch cfg.Cloud {
	case "aws":
		if strings.TrimSpace(cfg.AWS.Region) == "" {
			sl.ReportError(cfg.AWS.Region, "aws.region", "Region", "required_for_aws", "")
		}
	case "azure":
		if strings.TrimSpace(cfg.Azure.Location) == "" {
			sl.ReportError(cfg.Azure.Location, "azure.location", "Location", "required_for_azure", "")
		}
	case "gcp":
		if strings.TrimSpace(cfg.GCP.Region) == "" {
			sl.ReportError(cfg.GCP.Region, "gcp.region", "Region", "required_for_gcp", "")
		}
	}
}

func Validate(cfg Config) error {
	if err := validate.Struct(cfg); err != nil {
		return formatValidationErrors(err)
	}
	return nil
}

func formatValidationErrors(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	messages := make([]string, 0, len(validationErrors))
	for _, fe := range validationErrors {
		messages = append(messages, translateFieldError(fe))
	}
	return fmt.Errorf("config validation failed: %s", strings.Join(messages, "; "))
}

func translateFieldError(fe validator.FieldError) string {
	field := fe.Namespace()
	field = strings.TrimPrefix(field, "Config.")
	field = strings.ToLower(field)

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	case "required_for_aws":
		return "aws.region is required when cloud=aws"
	case "required_for_azure":
		return "azure.location is required when cloud=azure"
	case "required_for_gcp":
		return "gcp.region is required when cloud=gcp"
	case "required_if":
		return fmt.Sprintf("%s is required", field)
	default:
		return fmt.Sprintf("%s failed validation (%s)", field, fe.Tag())
	}
}
