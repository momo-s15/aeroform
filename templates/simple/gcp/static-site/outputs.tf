output "site_url" {
  description = "HTTP URL of the static site via the global forwarding rule"
  value       = "http://${google_compute_global_forwarding_rule.site.ip_address}"
}

output "bucket_name" {
  description = "GCS bucket hosting the static content"
  value       = google_storage_bucket.site.name
}
