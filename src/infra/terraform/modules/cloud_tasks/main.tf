resource "google_project_service" "cloud_tasks_api" {
  project = var.project_id
  service = "cloudtasks.googleapis.com"

  disable_on_destroy = false
}

resource "time_sleep" "wait_api_propagation" {
  depends_on      = [google_project_service.cloud_tasks_api]
  create_duration = "30s"
}

resource "google_cloud_tasks_queue" "fila_bacana" {
  name     = var.tasks_name
  location = var.region

  rate_limits {
    max_concurrent_dispatches = 50
    max_dispatches_per_second = 10.0
  }

  retry_config {
    max_attempts       = 5
    max_retry_duration = "3600s"
    min_backoff        = "10s"
    max_backoff        = "4000s"
    max_doublings      = 5
  }

  depends_on = [google_project_service.cloud_tasks_api]
}

resource "google_service_account" "cloud_task_sa" {
  project      = var.project_id
  account_id   = "cloud-task-editor"
  display_name = "Cloud Task Service Account"
}