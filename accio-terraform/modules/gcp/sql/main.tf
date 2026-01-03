variable "project_id" {
  type = string
}

variable "region" {
  type = string
}

variable "instance_name" {
  type = string
}

variable "database_version" {
  type    = string
  default = "POSTGRES_15"
}

resource "google_sql_database_instance" "this" {
  name             = var.instance_name
  database_version = var.database_version
  region           = var.region
  project          = var.project_id

  settings {
    tier = "db-f1-micro"
    user_labels = {
      managed_by = "accio"
    }
  }

  deletion_protection = false # for demo purposes
}

output "connection_name" {
  value = google_sql_database_instance.this.connection_name
}
