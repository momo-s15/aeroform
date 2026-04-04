terraform {
  required_version = ">= 1.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.80"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "site" {
  name     = "rg-${var.project_name}"
  location = var.location

  tags = {
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

resource "azurerm_storage_account" "site" {
  name                     = replace(var.project_name, "-", "")
  resource_group_name      = azurerm_resource_group.site.name
  location                 = azurerm_resource_group.site.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  min_tls_version          = "TLS1_2"

  blob_properties {
    versioning_enabled = true
  }

  tags = {
    Name      = var.project_name
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

# Static website config (nested static_website on storage account is deprecated in favour of this resource).
resource "azurerm_storage_account_static_website" "site" {
  storage_account_id = azurerm_storage_account.site.id

  index_document     = "index.html"
  error_404_document = "error.html"
}

# Classic Azure CDN (azurerm_cdn_profile / Standard_Microsoft) cannot be created after 2025-10-01.
# This template serves the site directly from Storage static website (HTTPS). Add Front Door in Pro Mode if you need a full CDN.
