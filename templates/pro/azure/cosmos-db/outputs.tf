output "account_id" {
  value = azurerm_cosmosdb_account.this.id
}

output "endpoint" {
  value = azurerm_cosmosdb_account.this.endpoint
}

output "primary_key" {
  value     = azurerm_cosmosdb_account.this.primary_key
  sensitive = true
}

output "database_name" {
  value = azurerm_cosmosdb_sql_database.this.name
}

output "connection_strings" {
  value     = azurerm_cosmosdb_account.this.connection_strings
  sensitive = true
}
