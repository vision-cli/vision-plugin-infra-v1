module "db" {
  source  = "GoogleCloudPlatform/sql-db/google//modules/postgresql"
  version = "15.1.0"

  name             = "${var.environment}-db"
  database_version = "POSTGRES_11"
  region           = var.region
  zone             = var.zone
  project_id       = module.project.project_id

  db_name           = var.db_name
  user_name         = var.db_user_name
  user_password     = var.db_user_password
  tier              = var.db_tier
  availability_type = "REGIONAL"
  disk_size         = var.db_vol_size

  #uncomment below to make DB deleteable for critical changes, consider this for prod
  #deletion_protection = false

  database_flags = [
    {
      name  = "cloudsql.iam_authentication",
      value = "on"
    }
  ]

  ip_configuration = {
    #disable public ip
    ipv4_enabled        = false
    private_network     = google_compute_network.private_network.id
    require_ssl         = false
    authorized_networks = []
    allocated_ip_range  = null
  }

# Dont enable backup for dev
#   backup_configuration = {
#     enabled                        = true
#     start_time                     = "00:00"
#     point_in_time_recovery_enabled = true
#     location                       = var.region
#     retained_backups               = 7
#     retention_unit                 = "COUNT"
#     transaction_log_retention_days = 6
#   }

  # try to avoid sqladmin service race-condition
  depends_on = [
    module.project,
    google_service_networking_connection.private_vpc_connection
  ]
}

# ------------------------------------------------------------
