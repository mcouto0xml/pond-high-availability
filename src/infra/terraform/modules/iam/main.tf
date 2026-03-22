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

# -----------------------------------------------------------------------------
# Permissão: escritor de Logs do Cloud Functions ServiceAccount
# -----------------------------------------------------------------------------
resource "google_project_iam_member" "function_sa_log_writer" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${var.function_sa_email}"
}