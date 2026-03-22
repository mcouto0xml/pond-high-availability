resource "google_project_service" "apis" {
  for_each = toset([
    "cloudfunctions.googleapis.com",   # Cloud Functions
    "cloudbuild.googleapis.com",       # Cloud Build (compila a função)
    "run.googleapis.com",              # Cloud Run (backend das Functions Gen 2)
    "storage.googleapis.com",          # GCS (upload do source code)
    "logging.googleapis.com",          # Logging
  ])

  project            = var.project_id
  service            = each.value
  disable_on_destroy = false
}

resource "google_service_account" "function_sa" {
  account_id   = "${var.function_name}-sa"
  display_name = "Service Account — ${var.function_name}"
}

resource "google_cloudfunctions2_function" "function" {
  name        = var.function_name
  location    = var.region
  description = "Função responsável por processar os dados "

  build_config {
    runtime     = "go125"                  # Ajuste para go121, go120 conforme necessário
    entry_point = "PostTelemetry" # Nome da função Go exportada

    source {
      storage_source {
        bucket = var.storage_name
        object = var.object_name
      }
    }
  }

  service_config {
    min_instance_count    = 0
    max_instance_count    = 50
    available_memory      = "256Mi"
    timeout_seconds       = 300
    service_account_email = google_service_account.function_sa.email
    
    environment_variables = var.environment_variables
  }

  depends_on = [google_project_service.apis]
}