output "app_url" {
  value = "https://${azurerm_linux_web_app.this.default_hostname}"
}

output "app_id" {
  value = azurerm_linux_web_app.this.id
}

output "identity_principal_id" {
  value = azurerm_linux_web_app.this.identity[0].principal_id
}

output "outbound_ip_addresses" {
  value = azurerm_linux_web_app.this.outbound_ip_addresses
}
