output "function_app_url" {
  description = "Default URL for the Azure Function App"
  value       = "https://${azurerm_linux_function_app.api.default_hostname}"
}

output "function_app_name" {
  description = "Name of the Function App"
  value       = azurerm_linux_function_app.api.name
}

output "resource_group" {
  description = "Resource group containing all resources"
  value       = azurerm_resource_group.api.name
}
