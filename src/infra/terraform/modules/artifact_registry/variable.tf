variable "project_id" {
    type = string
    description = "Projeto utilizado na ponderada"
}

variable "region" {
    type = string
    description = "Região utilizada para criar o Artifact Registry"
}

variable "repository_id" {
    type = string
    description = "Nome do repositório criado"
}

variable "description" {
    type = string
    description = "Descrição do repositório criado"
}

variable "enable_cleanup_policy" {
    type = bool
    description = "Habilitar política de limpeza"
}

variable "labels" {
  description = "Mapa de labels a serem aplicadas ao repositório."
  type        = map(string)
  default     = {}
}