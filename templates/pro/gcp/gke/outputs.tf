output "cluster_name" {
  value = google_container_cluster.this.name
}

output "cluster_endpoint" {
  value     = google_container_cluster.this.endpoint
  sensitive = true
}

output "cluster_ca_certificate" {
  value     = google_container_cluster.this.master_auth[0].cluster_ca_certificate
  sensitive = true
}

output "node_pool_name" {
  value = google_container_node_pool.default.name
}

output "workload_identity_pool" {
  value = google_container_cluster.this.workload_identity_config[0].workload_pool
}
