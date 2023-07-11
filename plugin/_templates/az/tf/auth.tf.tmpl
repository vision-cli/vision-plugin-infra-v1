## Random uuid for oauth permission scope
resource "random_uuid" "oauth2_perm_scope_user_read_uuid" {
}

## App registration
resource "azuread_application" "azconapp_auth_app_reg" {
  api {
    oauth2_permission_scope {
      admin_consent_description  = "Allows users to sign-in to the app, and allows the app to read the profile of signed-in users. It also allows the app to read basic company information of signed-in users."
      admin_consent_display_name = "Sign in and read user profile"
      enabled                    = true
      id                         = random_uuid.oauth2_perm_scope_user_read_uuid.result #Can generate random uuid -The unique identifier of the delegated permission. Must be a valid UUID
      type                       = "User"
      user_consent_description   = "Allows you to sign in to the app with your organizational account and let the app read your profile. It also allows the app to read basic company information."
      user_consent_display_name  = "Sign you in and read your profile"
      value                      = "User.Read" # value used for scp claim in OAuth 2.0 access tokens
    }
  }
  display_name            = "{{.AppName}}"
  group_membership_claims = ["All", "ApplicationGroup"] # Configures the groups claim issued in a user or OAuth 2.0 access token that the app expects

  identifier_uris               = ["api://${var.auth_identifier_uri}"] # A set of user-defined URI(s) that uniquely identify an application within its Azure AD tenant, or within a verified custom domain if the application is multi-tenant. The identifier_uri is mapped to the application id URI in the 'expose an API' section on the console. When setting it manually, the uuid value defaults to the application id. On terraform however, self-reference of a resource is not allowed as the application id is an attribute known after the terraform apply. It is possible to have a random value but make sure it is set to match the audience in the identity provider section.
  oauth2_post_response_required = false                                      # Specifies whether, as part of OAuth 2.0 token requests, Azure AD allows POST requests, as opposed to GET requests. Defaults to false, which specifies that only GET requests are allowed
  optional_claims {
    access_token {
      additional_properties = ["sam_account_name"]
      essential             = false
      name                  = "groups"
      # source =
    }
    id_token {
      additional_properties = ["sam_account_name"]
      essential             = false
      name                  = "groups"
      # source =

    }
    saml2_token {
      additional_properties = ["sam_account_name"]
      essential             = false
      name                  = "groups"
      # source =
    }
  }
  # owners = # By default, no owners are assigned. A set of object IDs of principals that will be granted ownership of the application. Supported object types are users or service principals. By default, no owners are assigned.
  prevent_duplicate_names = true # If true, will return an error if an existing application is found with the same name
  ## Values for the ids of the api permissions were found on Microsoft Docs: "https://learn.microsoft.com/en-gb/graph/permissions-reference#all-permissions-and-ids"
  required_resource_access {
    resource_app_id = "00000003-0000-0000-c000-000000000000" # Microsoft Graph Resource ID
    resource_access {
      id   = "e1fe6dd8-ba31-4d61-89e7-88639da4683d" # id of the User.Read api permission
      type = "Scope"

    }
    ## Commented out till there's a need by devs to display emails on sign-in
    # resource_access {
    #   id   = "64a6cdd6-aab1-4aaf-94b8-3cc8405e90d0" # id of the email api permission
    #   type = "Scope"

    # }
    # resource_access {
    #   id   = "37f7f235-527c-4136-accd-4a02d197296e" # id of the openid api permission
    #   type = "Scope"

    # }
    # resource_access {
    #   id   = "14dad69e-099b-42c9-810b-d002981feec1" # id of the profile api permission
    #   type = "Scope"
    # }

  }
  sign_in_audience = "AzureADMyOrg" # The Microsoft account types that are supported for the current application
  web {
    implicit_grant {
      access_token_issuance_enabled = false
      id_token_issuance_enabled     = true
    }
    redirect_uris = [
      "https://${var.app_aca_name}.${jsondecode(azapi_resource.containerapp_env.output).properties.defaultDomain}/.auth/login/aad/callback"
    ]
  }
}

## App registration Password for authentication
resource "azuread_application_password" "auth_app_registration_password" {
  application_object_id = azuread_application.azconapp_auth_app_reg.object_id
  display_name          = "az-containerapp-auth-secret"
  end_date_relative     = var.auth_app_reg_password_validity_period
}

data "azurerm_client_config" "current" {
}



## For full configuration of auth, the rest of the set up was done with az cli commands.
## A secret is set with the az containerapp secret set command - and this references the value placed in the keyvault - look at deployment pipelines for an example - tbd
## The auth config (created by the azapiprovider -lines 97-127 of the auth.tf file) was created without a clientSecretSettingName This will be  updated with the name of the secret set using the az containerapp auth microsoft update command so it can match the name of the secret created above.
