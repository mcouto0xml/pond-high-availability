output "registry_url" {
  description = "URL base para push/pull de imagens (Docker) ou pacotes."
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${var.repository_id}"
}

output "artifact_registry_sa_email" {
  description = "Email respectivo ao Service Account do Artifact Registry"
  value = google_service_account.artifact_registry_sa.email
}