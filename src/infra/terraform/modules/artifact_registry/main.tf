# Aqui usei o Claude para terraformar o artifact registry 

# ─────────────────────────────────────────────
# Habilita a API do Artifact Registry
# ─────────────────────────────────────────────
resource "google_project_service" "artifact_registry" {
  project            = var.project_id
  service            = "artifactregistry.googleapis.com"
  disable_on_destroy = false
}

# ─────────────────────────────────────────────
# Repositório do Artifact Registry
# ─────────────────────────────────────────────
resource "google_artifact_registry_repository" "this" {
  project       = var.project_id
  location      = var.region
  repository_id = var.repository_id
  description   = var.description
  format        = "DOCKER"

  dynamic "cleanup_policies" {
    for_each = var.enable_cleanup_policy ? [1] : []
    content {
      id     = "keep-minimum-versions"
      action = "KEEP"

      most_recent_versions {
        keep_count = 5
      }
    }
  }

  dynamic "cleanup_policies" {
    for_each = var.enable_cleanup_policy ? [1] : []
    content {
      id     = "delete-old-artifacts"
      action = "DELETE"

      condition {
        older_than = "2592000s"
      }
    }
  }

  labels = var.labels

  depends_on = [google_project_service.artifact_registry]
}

resource "google_service_account" "artifact_registry_sa" {
    project = var.project_id
    account_id = "artifact-registry-editor"
}
