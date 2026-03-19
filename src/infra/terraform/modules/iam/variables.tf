variable "project_id" {
  type        = string
  description = "Nome do projeto utilizado no GCP"
}

variable "region" {
  type        = string
  description = "Região utilizada no GCP"
}

variable "cloud_tasks_sa_email" {
  type        = string
  description = "Email utilizado pela Service Account do Cloud Tasks"
}

variable "admin_email" {
  type        = string
  description = "Email do Admin do Cloud Tasks"
}

variable "artifact_registry_sa_email" {
  type = string
  description = "Email respectivo a Service Account do Artifact Registry"
}

variable "ar_repository" {
  type = string
  description = "Nome do Repositório Artifact Registry"
}