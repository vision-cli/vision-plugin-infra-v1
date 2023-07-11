# Key vault secret for app registration password
resource "azurerm_key_vault_secret" "kv_auth_secret" {
  name         = var.keyvault_app_reg_password_name
  value        = azuread_application_password.auth_app_registration_password.value
  key_vault_id = data.azurerm_key_vault.keyvault.id
}

# Key vault secret for postgres password
resource "azurerm_key_vault_secret" "kv_db_password" {
  name         = var.keyvault_db_password_name
  value        = var.db_pass
  key_vault_id = data.azurerm_key_vault.keyvault.id
}

# Managed identity for use by the container app
resource "azurerm_user_assigned_identity" "aca_managed_identity" {
  name                = "aca-managed-identity"
  location            = var.location
  resource_group_name = var.rg_name
}

# Role assignment for the container app to pull secrets from the Key Vault
resource "azurerm_role_assignment" "aca_managed_id_role_assignment" {
  scope                = data.azurerm_resource_group.project_rg.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_user_assigned_identity.aca_managed_identity.principal_id
}

# https://ng-neptune-dev-kv.vault.azure.net/secrets/ng-neptune-aca-auth-reg-password-dev
