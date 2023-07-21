resource "google_compute_instance" "jumpbox" {
  name         = "jumpbox"
  machine_type = "e2-micro"
  zone         = var.zone
  project      = module.project.project_id

  tags = ["jumpbox-tag"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = google_compute_network.private_network.id

    access_config {
      // Ephemeral public IP
    }
  }

  service_account {
    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    email  = module.project.service_account_email
    scopes = ["cloud-platform"]
  }
}
