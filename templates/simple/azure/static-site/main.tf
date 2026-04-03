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

  static_website {
    index_document     = "index.html"
    error_404_document = "error.html"
  }

  blob_properties {
    versioning_enabled = true
  }

  tags = {
    Name      = var.project_name
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

resource "azurerm_cdn_profile" "site" {
  name                = "${var.project_name}-cdn"
  resource_group_name = azurerm_resource_group.site.name
  location            = "global"
  sku                 = "Standard_Microsoft"

  tags = {
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

resource "azurerm_cdn_endpoint" "site" {
  name                = var.project_name
  profile_name        = azurerm_cdn_profile.site.name
  resource_group_name = azurerm_resource_group.site.name
  location            = "global"

  origin_host_header = azurerm_storage_account.site.primary_web_host

  origin {
    name      = "storage"
    host_name = azurerm_storage_account.site.primary_web_host
  }

  delivery_rule {
    name  = "EnforceHTTPS"
    order = 1

    request_scheme_condition {
      operator     = "Equal"
      match_values = ["HTTP"]
    }

    url_redirect_action {
      redirect_type = "Found"
      protocol      = "Https"
    }
  }

  tags = {
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}
