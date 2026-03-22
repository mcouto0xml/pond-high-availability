variable "project_id" {
  type        = string
  description = "Projeto utilizado para a ponderada"
}

variable "region" {
  type        = string
  description = "Região utilizada no projeto"
}

variable "tasks_name" {
  type        = string
  description = "Nome do Cloud Tasks"
}

variable "supabase_access_token" {
  type        = string
  description = "Chave de acesso para o Supabase"

}

variable "database_region" {
  type        = string
  description = "Região utilizada pela Database"
}

variable "database_password" {
  type        = string
  description = "Senha da Database"
}

variable "supabase_org_id" {
  type        = string
  description = "Id da Organização do Supabase"
}

variable "admin_email" {
  type        = string
  description = "Email do usuário com permissão de Admin no Cloud Tasks"
}

variable "function_name" {
  type = string
  description = "Nome da Cloud Function"
}

variable "environment_variables" {
  type = map(string)
  description = "Variáveis de ambiente do Cloud Function"
}

variable "function_source_dir" {
  type = string
  description = "Diretório da Cloud Function"
}