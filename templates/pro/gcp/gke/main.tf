terraform {
  required_version = ">= 1.0"
}

resource "google_container_cluster" "this" {
  name     = var.name
  location = var.region
  project  = var.project_id

  initial_node_count       = 1
  remove_default_node_pool = true

  network    = var.network
  subnetwork = var.subnetwork

  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = var.private_endpoint
    master_ipv4_cidr_block  = var.master_cidr
  }

  ip_allocation_policy {}

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  release_channel {
    channel = "REGULAR"
  }

  master_auth {
    client_certificate_config {
      issue_client_certificate = false
    }
  }

  logging_service    = "logging.googleapis.com/kubernetes"
  monitoring_service = "monitoring.googleapis.com/kubernetes"
}

resource "google_container_node_pool" "default" {
  name     = "${var.name}-default"
  location = var.region
  project  = var.project_id
  cluster  = google_container_cluster.this.name

  initial_node_count = var.node_count

  autoscaling {
    min_node_count = var.min_node_count
    max_node_count = var.max_node_count
  }

  node_config {
    machine_type = var.machine_type
    disk_size_gb = 50
    disk_type    = "pd-standard"

    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]

    workload_metadata_config {
      mode = "GKE_METADATA"
    }

    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }
}
