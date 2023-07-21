output "load_balancer_ip" {
  value = module.http_lb.external_ip
}

output "service_account_email" {
  value = module.project.service_account_email
}

output "workload_identity_pool_id" {
  value = google_iam_workload_identity_pool.github_pool.id
}
