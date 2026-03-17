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

module "iam" {
  source = "./modules/iam"

  region               = var.region
  project_id           = var.project_id
  admin_email          = var.admin_email
  cloud_tasks_sa_email = module.cloud_tasks.cloud_task_sa_email
}