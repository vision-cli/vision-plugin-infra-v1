module "example" {
  source  = "GoogleCloudPlatform/cloud-run/google"
  version = "~> 0.2.0"

  # Required variables
  service_name           = "example"
  project_id             = var.project_id
  location               = var.region
  image                  = "gcr.io/cloudrun/hello"
  service_account_email  = var.service_acct_email
  service_annotations = {
    "run.googleapis.com/ingress"              : "internal-and-cloud-load-balancing"
    "run.googleapis.com/vpc-access-egress"    : "private-ranges-only"
    "run.googleapis.com/vpc-access-connector" : var.vpc_access_connector_name
  }
  members = [
    "allUsers", # external facing services only work with allUsers
  ]
}

resource "google_compute_region_network_endpoint_group" "example_neg" {
  name                  = "example-neg"
  network_endpoint_type = "SERVERLESS"
  project               = var.project_id
  region                = var.region
  cloud_run {
    service = module.example.service_name
  }
}
