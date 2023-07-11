
## Linux JB

## Public IP address for jumpbox
resource "azurerm_public_ip" "jumpbox_publicip" {
  name                = var.jb_public_ip
  location            = var.location
  resource_group_name = var.rg_name
  allocation_method   = "Dynamic"
}

## Network interface for jumpbox
resource "azurerm_network_interface" "jumpbox_nic" {
  name                = var.jb_nic
  location            = var.location
  resource_group_name = var.rg_name

  ip_configuration {
    name                          = var.jb_ip_conf
    subnet_id                     = azurerm_subnet.jumpbox_subnet.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.jumpbox_publicip.id
  }
}

# Create (and display) an SSH key
resource "tls_private_key" "lnx-jumpbox-sshkey" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

# Linux virtual machine to act as jumpbox
resource "azurerm_linux_virtual_machine" "jumpbox" {
  name                            = var.jb
  resource_group_name             = var.rg_name
  location                        = var.location
  size                            = "Standard_F2"
  admin_username                  = "azureuser"
  disable_password_authentication = true
  network_interface_ids = [
    azurerm_network_interface.jumpbox_nic.id
  ]

  os_disk {
    name                 = var.os_lnx_disk
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  admin_ssh_key {
    username   = "azureuser"
    public_key = tls_private_key.lnx-jumpbox-sshkey.public_key_openssh
  }
}

## Network Security group to allow SSH access only
resource "azurerm_network_security_group" "jumpbox_nsg" {
  name                = var.jb_nsg
  location            = var.location
  resource_group_name             = var.rg_name

  security_rule {
    name                       = "Allow-SSH"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

## Associate NSG with network interface for jumpbox
resource "azurerm_network_interface_security_group_association" "jumpbox_nsg_nic" {
  network_interface_id      = azurerm_network_interface.jumpbox_nic.id
  network_security_group_id = azurerm_network_security_group.jumpbox_nsg.id
}

# Auto shutdown virtual machine
resource "azurerm_dev_test_global_vm_shutdown_schedule" "jumpbox" {
  virtual_machine_id = azurerm_linux_virtual_machine.jumpbox.id
  location           = var.location
  enabled            = true

  daily_recurrence_time = "1900"
  timezone              = "GMT Standard Time"

  notification_settings {
    enabled = false
  }
}
