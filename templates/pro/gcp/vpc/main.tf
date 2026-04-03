terraform {
  required_version = ">= 1.0"
}

resource "google_compute_network" "this" {
  name                    = var.name
  auto_create_subnetworks = false
}
