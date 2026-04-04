# aeroform-schema: azure-static-site/3
terraform {
  required_version = ">= 1.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.80"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0"
    }
  }
}

# Storage account names are globally unique across Azure; suffix avoids collisions on common project slugs.
resource "random_string" "storage_suffix" {
  length  = 6
  special = false
  upper   = false
}

locals {
  # 3–24 chars, lowercase letters and numbers only (project_name allows only a-z, 0-9, hyphens — avoid regexreplace for older Terraform).
  name_slug            = substr(replace(lower(var.project_name), "-", ""), 0, 18)
  storage_account_name = "${local.name_slug}${random_string.storage_suffix.result}"
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
  name                     = local.storage_account_name
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
