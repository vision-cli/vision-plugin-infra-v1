## Create a file share in Storage account
resource "azurerm_storage_share" "strg_share" {
  name                 = var.strg_fileshare_name
  storage_account_name = var.strg_acc_name
  quota                = 5
}

## Create a file share in the container app env to connect to file share storage account
resource "azurerm_container_app_environment_storage" "aca_env_strg" {
  name                         = var.strg_aca_env_strg_name
  container_app_environment_id = azapi_resource.containerapp_env.id
  account_name                 = var.strg_acc_name
  share_name                   = azurerm_storage_share.strg_share.name
  access_key                   = data.azurerm_storage_account.strg_acc.primary_access_key
  access_mode                  = "ReadWrite"
}

## Database DNS zone
resource "azurerm_private_dns_zone" "priv_dns" {
  name                = var.db_priv_dns_zone
  resource_group_name = var.rg_name
}

## Database DNS zone link to main vnet
resource "azurerm_private_dns_zone_virtual_network_link" "db_dns_vnetlink" {
  name                  = var.db_dns_vnetlink_name
  resource_group_name   = var.rg_name
  private_dns_zone_name = azurerm_private_dns_zone.priv_dns.name
  virtual_network_id    = azurerm_virtual_network.project_vnet.id
}

## Postgresql flexible server
resource "azurerm_postgresql_flexible_server" "pg_db" {
  name                   = var.db_server_name
  resource_group_name    = var.rg_name
  location               = var.location
  version                = "14"
  delegated_subnet_id    = azurerm_subnet.db_subnet.id
  private_dns_zone_id    = azurerm_private_dns_zone.priv_dns.id
  administrator_login    = "pgsqladmin"
  administrator_password = var.db_pass
  zone                   = "1"

  storage_mb = 32768

  sku_name = "B_Standard_B1ms" #"GP_Standard_D2s_v3"

  depends_on = [
    azurerm_private_dns_zone_virtual_network_link.db_dns_vnetlink
  ]
}
