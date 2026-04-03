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

resource "azurerm_resource_group" "api" {
  name     = "rg-${var.project_name}"
  location = var.location

  tags = {
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

resource "azurerm_storage_account" "api" {
  name                     = replace(var.project_name, "-", "")
  resource_group_name      = azurerm_resource_group.api.name
  location                 = azurerm_resource_group.api.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  min_tls_version          = "TLS1_2"

  tags = {
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

resource "azurerm_service_plan" "api" {
  name                = "${var.project_name}-plan"
  resource_group_name = azurerm_resource_group.api.name
  location            = azurerm_resource_group.api.location
  os_type             = "Linux"
  sku_name            = "Y1"

  tags = {
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}

resource "azurerm_linux_function_app" "api" {
  name                = var.project_name
  resource_group_name = azurerm_resource_group.api.name
  location            = azurerm_resource_group.api.location

  storage_account_name       = azurerm_storage_account.api.name
  storage_account_access_key = azurerm_storage_account.api.primary_access_key
  service_plan_id            = azurerm_service_plan.api.id

  https_only = true

  site_config {
    application_stack {
      node_version = "20"
    }

    ftps_state = "Disabled"
  }

  app_settings = {
    "FUNCTIONS_WORKER_RUNTIME" = "node"
  }

  tags = {
    ManagedBy = "aeroform"
    AeroMode  = "simple"
  }
}
