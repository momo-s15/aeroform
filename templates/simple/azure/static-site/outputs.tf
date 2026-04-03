output "cdn_endpoint_url" {
  description = "CDN endpoint URL for the static site"
  value       = "https://${azurerm_cdn_endpoint.site.fqdn}"
}

output "storage_account_name" {
  description = "Storage account hosting the static content"
  value       = azurerm_storage_account.site.name
}

output "resource_group" {
  description = "Resource group containing all resources"
  value       = azurerm_resource_group.site.name
}
