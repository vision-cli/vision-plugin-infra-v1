module "project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 14.2"

  name                 = "${var.project_name}-${var.environment}-${var.unique_str}"
  random_project_id    = false
  org_id               = var.org_id
  billing_account      = var.billing_account
  folder_id            = var.folder_id

  activate_apis = [
    "compute.googleapis.com",
    "iam.googleapis.com",
    "iap.googleapis.com",
    "run.googleapis.com",
    "vpcaccess.googleapis.com",
    "servicenetworking.googleapis.com",
  ]
}

resource "google_project_service_identity" "cloudrun_sa" {
  provider = google-beta

  project = module.project.project_id
  service = "iap.googleapis.com"
}
