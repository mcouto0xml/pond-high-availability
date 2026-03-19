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

variable "ar_repository_id" {
  type = string
  description = "Nome do Artifact Registry"
} 

variable "ar_description" {
  type = string
  description = "Descrição do Artifact Registry"
}

variable "ar_enable_cleanup_policy" {
  type = bool
  description = "Política de limpeza"
}

variable "ar_labels" {
  type = map(string)
  description = "Labels do Artifact Registry"
}