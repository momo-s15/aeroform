terraform {
  required_version = ">= 1.0"
}

resource "google_sql_database_instance" "this" {
  name             = var.name
  project          = var.project_id
  region           = var.region
  database_version = var.database_version

  deletion_protection = false

  settings {
    tier              = var.tier
    availability_type = var.high_availability ? "REGIONAL" : "ZONAL"
    disk_size         = var.disk_size
    disk_type         = "PD_SSD"
    disk_autoresize   = true

    ip_configuration {
      ipv4_enabled    = false
      private_network = var.network
      require_ssl     = true
    }

    backup_configuration {
      enabled                        = true
      start_time                     = "03:00"
      point_in_time_recovery_enabled = var.database_version != "MYSQL_8_0"
      transaction_log_retention_days = 7

      backup_retention_settings {
        retained_backups = 7
      }
    }

    maintenance_window {
      day          = 7
      hour         = 4
      update_track = "stable"
    }

    database_flags {
      name  = "log_checkpoints"
      value = "on"
    }
  }
}

resource "google_sql_database" "default" {
  name     = var.db_name
  project  = var.project_id
  instance = google_sql_database_instance.this.name
}

resource "google_sql_user" "default" {
  name     = var.db_user
  project  = var.project_id
  instance = google_sql_database_instance.this.name
  password = var.db_password
}
