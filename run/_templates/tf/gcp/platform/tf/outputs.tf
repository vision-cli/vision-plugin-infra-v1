output "load_balancer_ip" {
  value = module.http_lb.external_ip
}

output "service_account_email" {
  value = module.project.service_account_email
}

output "workload_identity_pool_id" {
  value = module.wif.pool_id
}

# ------------------------------------------------------------
