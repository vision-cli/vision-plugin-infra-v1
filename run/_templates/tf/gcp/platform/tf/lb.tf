module "http_lb" {
  source            = "GoogleCloudPlatform/lb-http/google//modules/serverless_negs"
  version           = "~> 9.0"

  project           = module.project.project_id
  name              = "${var.environment}-http-lb"

  ssl                             = true
  managed_ssl_certificate_domains = [var.domain]
  https_redirect                  = true
  create_url_map                  = false
  url_map                         = module.cloud_run.urlmap

  backends = module.cloud_run.backends
}


# ------------------------------------------------------------


