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

module "wif" {
  source     = "SudharsaneSivamany/workload-identity-federation-multi-provider/google"

  project_id = module.project.project_id
  pool_id    = "github-pool"
  wif_providers = [
  { provider_id          = "github-provider"
    select_provider      = "oidc"
    provider_config      = {
                             issuer_uri = "https://token.actions.githubusercontent.com"
                             allowed_audiences = "https://iam.googleapis.com/projects/${module.project.project_number}/locations/global/workloadIdentityPools/github-pool/providers/github-provider"
                           }
    disabled             = false
    attribute_mapping    = {
                             "attribute.actor"      = "assertion.actor"
                             "attribute.repository" = "assertion.repository"
                             "google.subject"       = "assertion.sub"
                           }
  },
]
  service_accounts = [
    {
      name           = module.project.service_account_email
#      attribute      = "attribute.repository/my-org/my-repo"
      all_identities = true
      roles          = [
        "roles/iam.serviceAccountUser",
        "roles/iam.serviceAccountOpenIdTokenCreator",
        "roles/iam.workloadIdentityUser",
        ]
    }
  ]
}
