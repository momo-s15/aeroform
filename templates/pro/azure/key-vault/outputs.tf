output "vault_id" {
  value = azurerm_key_vault.this.id
}

output "vault_uri" {
  value = azurerm_key_vault.this.vault_uri
}

output "vault_name" {
  value = azurerm_key_vault.this.name
}

output "tenant_id" {
  value = azurerm_key_vault.this.tenant_id
}
