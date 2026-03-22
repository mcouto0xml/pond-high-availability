resource "google_project_service" "apis" {
  for_each = toset([
    "cloudtasks.googleapis.com",
  ])

  project            = var.project_id
  service            = each.value
  disable_on_destroy = false
}

resource "google_storage_bucket" "function_source" {
  name                        = "${var.project_id}-function-source"
  location                    = var.region
  force_destroy               = true
  uniform_bucket_level_access = true

  # Evita que versões antigas do source acumulem indefinidamente
  lifecycle_rule {
    condition { num_newer_versions = 3 }
    action    { type = "Delete" }
  }

  depends_on = [google_project_service.apis]
}

data "archive_file" "function_zip" {
  type        = "zip"
  source_dir  = var.function_source_dir
  output_path = "${path.module}/.tmp/function_source.zip"

  # Ignora arquivos desnecessários no zip
  excludes = [
    ".git",
    ".gitignore",
    "*.test",
  ]
}

resource "google_storage_bucket_object" "function_source" {
  name   = "source-${data.archive_file.function_zip.output_md5}.zip"
  bucket = google_storage_bucket.function_source.name
  source = data.archive_file.function_zip.output_path
}