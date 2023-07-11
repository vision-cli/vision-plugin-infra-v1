### Base infrastructure ###
/*************************/

## Resource group
#  Pre-created for in order to assign AD privileges for developers
#  See variables.tf and ./config/poc.tfvars

## Log analytics
resource "azurerm_log_analytics_workspace" "loganalytics" {
  name                = "${var.app_aca_name}-log"
  location            = var.location
  resource_group_name = var.rg_name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

### Container applicaton ###
/**************************/

## Container app environment

resource "azapi_resource" "containerapp_env" {
  type      = "Microsoft.App/managedEnvironments@2022-03-01" # 2022-10-01
  name      = "${var.app_aca_name}-env"
  parent_id = data.azurerm_resource_group.project_rg.id
  location  = var.location
  body = jsonencode({
    properties = {
      appLogsConfiguration = {
        destination = "log-analytics"
        logAnalyticsConfiguration = {
          customerId = azurerm_log_analytics_workspace.loganalytics.workspace_id
          sharedKey  = azurerm_log_analytics_workspace.loganalytics.primary_shared_key
        } ## This section is required by terraform as the log analytics workspace keeps a record of the system logs and console logs for container apps by default
      }
      vnetConfiguration = {
        internal               = false # needs to be false to ensure container app env. is exteranlly accessible
        infrastructureSubnetId = azurerm_subnet.app_subnet.id
        dockerBridgeCidr       = var.dockerbridgeCidr
        platformReservedCidr   = var.platformReservedCidr
        platformReservedDnsIP  = var.platformReservednsIP
      }
    }
  })

  response_export_values  = ["properties.defaultDomain", "properties.staticIp"]
  ignore_missing_property = true

  depends_on = [
    azurerm_virtual_network.project_vnet
  ]
}

resource "azapi_resource" "graphql-server" {
  type      = "Microsoft.App/containerapps@2022-03-01"
  name      = "${var.app_aca_name}-graphql-server"
  parent_id = data.azurerm_resource_group.project_rg.id
  location  = var.location
  body = jsonencode({
    properties = {
      managedEnvironmentId = azapi_resource.containerapp_env.id
      configuration = {
        ingress = {
          external : true,
          targetPort : 8080,
          # customDomains = [
          #   {
          #     name          = var.go_domain_name,
          #     certificateId = data.azurerm_container_app_environment_certificate.ca_env_cert.id,
          #     bindingType   = "SniEnabled"
          #   }
          # ]
        },

        registries = [
          {
            server            = data.azurerm_container_registry.acr.login_server
            username          = data.azurerm_container_registry.acr.admin_username
            passwordSecretRef = "registry-password"
          }
        ],
        secrets : [
          {
            name = "registry-password"
            # Todo: Container apps does not yet support Managed Identity connection to ACR
            value = data.azurerm_container_registry.acr.admin_password
          }
        ]
      },

      template = {
        containers = [
          {
            image = "${data.azurerm_container_registry.acr.login_server}/graphql/server"
            name  = "graphql-server"
            env   = [
              {
                name: "DATABASE_URL",
                value: "postgres://pgsqladmin:${var.db_pass}@${azurerm_postgresql_flexible_server.pg_db.fqdn}/postgres?sslmode=require"
              },
              {
                name: "GRAPHIQL_ENABLED",
                value: "true"
              }
            ]
            resources = {
              cpu    = 0.5
              memory = "1.0Gi"
            }
          }
        ]
        scale = {
          minReplicas = 1,
          maxReplicas = 3,
          rules = [
            {
              name = "connections",
              http = {
                metadata = {
                  concurrentRequests = "25"
                }
              }
            }
          ]
        }
      }
    }

  })
  depends_on = [
    azapi_resource.containerapp_env
  ]
  lifecycle {
    ignore_changes = [
      body
    ]
  }
}

