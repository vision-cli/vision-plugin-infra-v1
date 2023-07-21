data "google_iam_policy" "users" {
  binding {
    role    = "roles/iap.httpsResourceAccessor"
    members = var.members
  }
}

resource "google_iap_web_backend_service_iam_policy" "backend_access_policy" {
  for_each = module.cloud_run.backends
  project             = module.project.project_id
  web_backend_service = "${var.environment}-http-lb-backend-${each.key}"
  policy_data         = data.google_iam_policy.users.policy_data
  depends_on = [
    module.http_lb
  ]
}

resource "google_iam_workload_identity_pool" "github_pool" {
  project                   = module.project.project_id
  workload_identity_pool_id = "github-pool"
}

resource "google_iam_workload_identity_pool_provider" "github_provider" {
  project                            = module.project.project_id
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-provider"
  display_name                       = "Github provider"
  description                        = "OIDC identity pool provider for cicd"
  disabled                           = false
  attribute_mapping                  = {
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
    "google.subject"       = "assertion.sub"
  }
  oidc {
    allowed_audiences = ["https://iam.googleapis.com/projects/${module.project.project_number}/locations/global/workloadIdentityPools/github-pool/providers/github-provider"]
    issuer_uri        = "https://token.actions.githubusercontent.com"
  }
}

