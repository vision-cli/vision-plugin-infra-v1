output "urlmap" {
  value = google_compute_url_map.urlmap.self_link
}

output "backends" {
  value = {
    example = {
      description                     = null
      protocol                        = "HTTP"
      port_name                       = "http"
      enable_cdn                      = false
      custom_request_headers          = null
      custom_response_headers         = null
      security_policy                 = null
      compression_mode                = null


      log_config = {
        enable = true
        sample_rate = 1.0
      }

      groups = [
        {
            group = google_compute_region_network_endpoint_group.example_neg.id
        },
      ]

      iap_config = {
        enable               = true
        oauth2_client_id     = var.oauth2_client_id
        oauth2_client_secret = var.oauth2_client_secret
      }
    }

  }
}