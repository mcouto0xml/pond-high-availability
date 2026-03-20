# Aqui utilizei o Claude para gerar as permissões
# Obs.: Precisaria saber as Roles de cabeça, coisa que não sei
# Aqui to utilizando o princípio da permissão máxima KKKKKKKKKKK
# To dando W/R para uma SA só, mas sei que em Prod precisaria separar isso serviço a serviço

# -----------------------------------------------------------------------------
# Permissão: enqueue (push) de tarefas no Cloud Tasks
# -----------------------------------------------------------------------------

resource "google_project_iam_member" "cloud_task_enqueuer" {
  project = var.project_id
  role    = "roles/cloudtasks.enqueuer"
  member  = "serviceAccount:${var.cloud_tasks_sa_email}"
}

# -----------------------------------------------------------------------------
# Permissão: dequeue/pull de tarefas no Cloud Tasks
# -----------------------------------------------------------------------------

resource "google_project_iam_member" "cloud_task_viewer" {
  project = var.project_id
  role    = "roles/cloudtasks.viewer"
  member  = "serviceAccount:${var.cloud_tasks_sa_email}"
}

resource "google_project_iam_member" "cloud_task_invoker" {
  project = var.project_id
  role = "roles/run.invoker"
  member = "serviceAccount:${var.cloud_tasks_sa_email}"
}

# -----------------------------------------------------------------------------
# Permissão: execução das tasks (invocar o target — Cloud Run / HTTP)
# -----------------------------------------------------------------------------

resource "google_project_iam_member" "cloud_task_runner" {
  project = var.project_id
  role    = "roles/cloudtasks.taskRunner"
  member  = "serviceAccount:${var.cloud_tasks_sa_email}"
}

# -----------------------------------------------------------------------------
# Permissão: admin completo no Cloud Tasks
# -----------------------------------------------------------------------------

resource "google_project_iam_member" "admin_cloud_tasks" {
  project = var.project_id
  role    = "roles/cloudtasks.admin"
  member  = "user:${var.admin_email}"
}

# -----------------------------------------------------------------------------
# Permissão: visualizar recursos do projeto (logs, monitoring)
# -----------------------------------------------------------------------------

resource "google_project_iam_member" "admin_viewer" {
  project = var.project_id
  role    = "roles/viewer"
  member  = "user:${var.admin_email}"
}

resource "google_artifact_registry_repository_iam_member" "ar_admin" {
  project = var.project_id
  location = var.region
  repository = var.ar_repository
  role = "roles/artifactregistry.admin"
  member = "user:${var.admin_email}"
}

# ─────────────────────────────────────────────
# IAM – Leitores do repositório Artifact Registry
# ─────────────────────────────────────────────
resource "google_artifact_registry_repository_iam_member" "reader" {
  project    = var.project_id
  location   = var.region
  repository = var.ar_repository
  role       = "roles/artifactregistry.reader"
  member     = "serviceAccount:${var.artifact_registry_sa_email}"
}

# ─────────────────────────────────────────────
# IAM – Escritores do repositório Artifact Registry
# ─────────────────────────────────────────────
resource "google_artifact_registry_repository_iam_member" "writer" {
  project    = var.project_id
  location   = var.region
  repository = var.ar_repository
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:${var.artifact_registry_sa_email}"
}
