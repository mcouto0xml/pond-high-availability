variable "project_id" {
  type        = string
  description = "Projeto utilizado para a ponderada"
}

variable "database_password" {
  type        = string
  description = "Senha para o Banco de Dados"
}

variable "database_region" {
  type        = string
  description = "Região utilizada no Banco de Dados"
}

variable "supabase_org_id" {
  type        = string
  description = "Id da organização utilizada no banco de dados"
}

variable "supabase_access_token" {
  type = string
  description = "Acess Token do Supabase"
}