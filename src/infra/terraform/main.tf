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

module "cloud_storage" {
  source = "./modules/cloud_storage"

  project_id = var.project_id
  region = var.region
  function_source_dir = var.function_source_dir
}

module "cloud_functions" {
  source = "./modules/cloud_function"

  project_id = var.project_id
  region = var.region
  function_name = var.function_name
  environment_variables = var.environment_variables
  storage_name = module.cloud_storage.storage_name
  object_name = module.cloud_storage.object_name
}

module "iam" {
  source = "./modules/iam"

  region               = var.region
  project_id           = var.project_id
  admin_email          = var.admin_email
  cloud_tasks_sa_email = module.cloud_tasks.cloud_task_sa_email
  function_sa_email    = module.cloud_functions.function_sa_email
}