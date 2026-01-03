variable "project_id" {
  type = string
}

variable "location" {
  type = string
}

variable "bucket_name" {
  type = string
}

resource "google_storage_bucket" "this" {
  name          = var.bucket_name
  location      = var.location
  project       = var.project_id
  force_destroy = true

  uniform_bucket_level_access = true

  labels = {
    managed_by = "accio"
  }
}

output "bucket_url" {
  value = google_storage_bucket.this.url
}
