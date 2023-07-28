# Create a main vnet
resource "azurerm_virtual_network" "project_vnet" {
  name                = var.vnet_name
  location            = var.location
  resource_group_name = var.rg_name
  address_space       = [var.vnet_cidr_block]
}

# Subnet for container app
resource "azurerm_subnet" "app_subnet" {
  name                 = var.vnet_app_subnet
  resource_group_name  = var.rg_name
  virtual_network_name = azurerm_virtual_network.project_vnet.name
  address_prefixes     = [var.vnet_app_cidr_block]
}

## Subnet for db
resource "azurerm_subnet" "db_subnet" {
  name                 = var.vnet_db_subnet
  resource_group_name  = var.rg_name
  virtual_network_name = azurerm_virtual_network.project_vnet.name
  address_prefixes     = [var.vnet_db_cidr_block]
  #   service_endpoints    = ["Microsoft.Storage"] This service endpoint is automatically added

  delegation {
    name = "delegation"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action"
      ]
    }
  }
  depends_on = [
    azurerm_virtual_network.project_vnet
  ]
}

# Subnet for Jumpbox
resource "azurerm_subnet" "jumpbox_subnet" {
  name                 = var.vnet_jumpbox_subnet
  resource_group_name  = var.rg_name
  virtual_network_name = azurerm_virtual_network.project_vnet.name
  address_prefixes     = [var.vnet_jumpbox_cidr_block]

  depends_on = [
    azurerm_virtual_network.project_vnet
  ]
}

# ## Ng-Neptune and Common Bastion Peering
# resource "azurerm_virtual_network_peering" "project_to_common_bastion_peering" {
#   name                         = var.project_common_bastion_peering
#   resource_group_name          = var.rg_name
#   virtual_network_name         = azurerm_virtual_network.project_vnet.name
#   remote_virtual_network_id    = data.azurerm_virtual_network.bastion_vnet.id
#   allow_virtual_network_access = true
#   allow_forwarded_traffic      = true
#   use_remote_gateways          = false
# }

# resource "azurerm_virtual_network_peering" "common_bastion_to_project_peering" {
#   name                         = var.common_bastion_project_peering
#   resource_group_name          = data.azurerm_virtual_network.bastion_vnet.resource_group_name
#   virtual_network_name         = data.azurerm_virtual_network.bastion_vnet.name
#   remote_virtual_network_id    = azurerm_virtual_network.project_vnet.id
#   allow_virtual_network_access = true
#   allow_forwarded_traffic      = true
#   use_remote_gateways          = false
# }
