output "website_url" {
  description = "HTTPS URL for the storage static website"
  value       = azurerm_storage_account.site.primary_web_endpoint
}

output "storage_account_name" {
  description = "Storage account hosting the static content"
  value       = azurerm_storage_account.site.name
}

output "resource_group" {
  description = "Resource group containing all resources"
  value       = azurerm_resource_group.site.name
}

output "upload_command" {
  description = "Upload local site files to the $web container (use your path instead of ./my-site)"
  value       = "az storage blob upload-batch --account-name ${azurerm_storage_account.site.name} -d '$web' -s ./my-site --auth-mode login --overwrite"
}
