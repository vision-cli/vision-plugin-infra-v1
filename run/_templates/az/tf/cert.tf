### Custom DNS & certificate binding ###
/**************************************/

## Create custom FQDN in Azure DNS

## Azure Container app Environment Default Domain
## Switched to custom domain blocks within the container app as the CNAME record here points to the container app environment and not the individual apps.
# resource "azurerm_dns_cname_record" "containerapp_cname_record" {
#   name                = var.atrias_poc_dns_cname_record_url
#   zone_name           = data.azurerm_dns_zone.atos_cerebro_common_dns_zone.name
#   resource_group_name = data.azurerm_dns_zone.atos_cerebro_common_dns_zone.resource_group_name
#   ttl                 = 300
#   record              = jsondecode(azapi_resource.atrias_containerapp_env.output).properties.defaultDomain
# }

## Container app environment certificate manually uploaded for an extra security layer.
#data "azurerm_container_app_environment_certificate" "ca_env_cert" {
#  name                         = "atos-test-cert"
#  container_app_environment_id = azapi_resource.containerapp_env.id
#}

#resource "azurerm_container_app_environment_certificate" "ca_env_cert" {
#  name                         = "project-cert"
#  container_app_environment_id = azapi_resource.containerapp_env.id
#  certificate_blob_base64      = filebase64("./cert.pem")
#  certificate_password         = ""
#}
