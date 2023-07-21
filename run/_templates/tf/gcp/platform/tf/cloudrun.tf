module "cloud_run" {
  source                    = "./cloud-run"
  project_id                = module.project.project_id
  region                    = var.region
  service_acct_email        = module.project.service_account_email
  oauth2_client_id          = var.oauth2_client_id
  oauth2_client_secret      = var.oauth2_client_secret
  members                   = var.members
  vpc_access_connector_name = google_vpc_access_connector.vpc_connector.name
  domain                    = var.domain
}
