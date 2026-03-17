output "project_id" {
  value       = supabase_project.ponderada.id
  description = "ID do projeto Supabase"
}

output "database_url" {
  value       = "postgresql://postgres:${var.database_password}@db.${supabase_project.ponderada.id}.supabase.co:5432/postgres"
  description = "Connection string do banco"
  sensitive   = true
}

output "api_url" {
  value       = "https://${supabase_project.ponderada.id}.supabase.co"
  description = "URL da API REST do projeto"
}