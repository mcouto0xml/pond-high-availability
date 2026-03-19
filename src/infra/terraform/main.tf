module "cloud_tasks" {
  source = "./modules/cloud_tasks"

  region     = var.region
  tasks_name = var.tasks_name
  project_id = var.project_id
}

module "postgresql" {
  source = "./modules/postgresql"

  database_region   = var.database_region
  database_password = var.database_password
  project_id        = var.project_id
  supabase_org_id   = var.supabase_org_id
  supabase_access_token = var.supabase_access_token
}

module "artifact_registry" {
  source = "./modules/artifact_registry"

  project_id = var.project_id
  region = var.region
  repository_id = var.ar_repository_id
  description = var.ar_description
  enable_cleanup_policy = var.ar_enable_cleanup_policy
  labels = var.ar_labels
}

module "iam" {
  source = "./modules/iam"

  region               = var.region
  project_id           = var.project_id
  admin_email          = var.admin_email
  cloud_tasks_sa_email = module.cloud_tasks.cloud_task_sa_email
  artifact_registry_sa_email = module.artifact_registry.artifact_registry_sa_email
  ar_repository = var.ar_repository_id
}