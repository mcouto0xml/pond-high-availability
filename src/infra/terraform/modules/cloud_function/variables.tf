variable "project_id" {
    type = string
    description = "ID do Projeto utilizado na Ponderada"
}

variable "function_name" {
    type = string
    description = "Nome da Cloud Function"
}

variable "region" {
    type = string
    description = "Região do Cloud Function"
}

variable "environment_variables" {
    type = map(string)
    description = "Dicionário das variáveis de ambiente"
}

variable "storage_name" {
  type = string
  description = "Nome do Cloud Storage"
}

variable "object_name" {
  type = string
  description = "Nome do objeto criado no Cloud Storage"
}