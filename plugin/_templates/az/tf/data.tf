## Project resource group created by Vision
data "azurerm_resource_group" "project_rg" {
  name = var.rg_name
}

## Storage account created by Vision
data "azurerm_storage_account" "strg_acc" {
  name = var.strg_acc_name
  resource_group_name = var.rg_name
}

## Keyvault created by Vision
data "azurerm_key_vault" "keyvault" {
  name = var.keyvault_name
  resource_group_name = var.rg_name
}

## Container registry created by Vision
data "azurerm_container_registry" "acr" {
  name = var.acr_name
  resource_group_name = var.rg_name
}

## Common bastion vnet
data "azurerm_virtual_network" "bastion_vnet" {
  name                = var.common_bastion_vnet
  resource_group_name = var.common_bastion_rg
}

## Existing Azure Dns Zone
data "azurerm_dns_zone" "atos_cerebro_common_dns_zone" {
  name                = var.atos_cerebro_common_dns_zone
  resource_group_name = var.atos_cerebro_common_dns_resource_group
}

resource "azurerm_dns_cname_record" "containerapp_cname_record" {
  name                = var.project_dns_cname_record_url
  zone_name           = data.azurerm_dns_zone.atos_cerebro_common_dns_zone.name
  resource_group_name = data.azurerm_dns_zone.atos_cerebro_common_dns_zone.resource_group_name
  ttl                 = 300
  record              = jsondecode(azapi_resource.containerapp_env.output).properties.defaultDomain
}

