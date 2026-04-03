terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_storage_bucket" "site" {
  name     = "${var.project_name}-site"
  project  = var.project_id
  location = var.region

  uniform_bucket_level_access = true
  force_destroy               = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "error.html"
  }

  versioning {
    enabled = true
  }

  labels = {
    managed-by = "aeroform"
    aero-mode  = "simple"
  }
}

resource "google_storage_bucket_iam_member" "public_read" {
  bucket = google_storage_bucket.site.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

resource "google_compute_backend_bucket" "site" {
  name        = "${var.project_name}-backend"
  bucket_name = google_storage_bucket.site.name
  enable_cdn  = true
}

resource "google_compute_url_map" "site" {
  name            = "${var.project_name}-urlmap"
  default_service = google_compute_backend_bucket.site.id
}

resource "google_compute_target_http_proxy" "site" {
  name    = "${var.project_name}-http-proxy"
  url_map = google_compute_url_map.site.id
}

resource "google_compute_global_forwarding_rule" "site" {
  name       = "${var.project_name}-fwd"
  target     = google_compute_target_http_proxy.site.id
  port_range = "80"
}
