resource "google_compute_url_map" "urlmap" {
  project     = var.project_id
  name        = "urlmap"
  description = "the url mapping to backends"

  default_service = "https://www.googleapis.com/compute/v1/projects/${var.project_id}/global/backendServices/dev-http-lb-backend-example"

  host_rule {
    hosts        = ["${var.domain}"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = "https://www.googleapis.com/compute/v1/projects/${var.project_id}/global/backendServices/dev-http-lb-backend-example"

    path_rule {
      paths   = ["/example/*"]
      service = "https://www.googleapis.com/compute/v1/projects/${var.project_id}/global/backendServices/dev-http-lb-backend-example"
    }
  }

  depends_on = [ 
    google_compute_region_network_endpoint_group.example_neg, 
  ]
}
