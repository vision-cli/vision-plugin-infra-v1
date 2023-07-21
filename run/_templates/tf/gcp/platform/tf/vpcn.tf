resource "google_compute_network" "private_network" {
  project = module.project.project_id
  name    = "private-network"
}

resource "google_compute_global_address" "private_ip_address" {
  project       = module.project.project_id
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.private_network.id
}

#Connector for use by app engine to privately addressed resources
resource "google_vpc_access_connector" "vpc_connector" {
  region        = var.region
  project       = module.project.project_id
  name          = "private-vpc-connector"
  ip_cidr_range = "10.8.0.0/28"
  network       = google_compute_network.private_network.name
  depends_on = [
    module.project
  ]
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.private_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

resource "google_compute_firewall" "allow-ssh-from-iap" {
  name        = "allow-ssh-from-iap"
  network     = google_compute_network.private_network.id
  project     = module.project.project_id
  target_tags = ["jumpbox-tag"]

  source_ranges = [
    "0.0.0.0/0",
  ]

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  depends_on = [
    # The project's services must be set up before the
    # network is enabled as the compute API will not
    # be enabled and cause the setup to fail.
    module.project,
  ]
}

