# Aqui utilizei o Claude para Terraformar um Banco de Dados no Supabase para mim

terraform {
  required_providers {
    supabase = {
      source  = "supabase/supabase"
      version = "~> 1.0"
    }
  }
}

# -----------------------------------------------------------------------------
# Projeto Supabase
# -----------------------------------------------------------------------------

resource "supabase_project" "ponderada" {
  organization_id   = var.supabase_org_id
  name              = var.project_id
  database_password = var.database_password
  region            = var.database_region

  # Aguarda o projeto ficar pronto antes de continuar
  lifecycle {
    ignore_changes = [database_password]
  }
}

# -----------------------------------------------------------------------------
# Schema — tabela telemetry
# Aplicado via SQL usando o provider de settings do Supabase
# -----------------------------------------------------------------------------

resource "supabase_settings" "db_schema" {
  project_ref = supabase_project.ponderada.id

  # O schema é definido no arquivo schema.json e executado como SQL
  # Para aplicar o DDL, utilize o bloco de migração abaixo via API REST do Supabase.
  # Em pipelines CI/CD, prefira usar `supabase db push` com as migrations.
}

resource "null_resource" "apply_schema" {
  depends_on = [supabase_project.ponderada]

  triggers = {
    schema_hash = filemd5("${path.module}/schema.json")
  }

  provisioner "local-exec" {
    environment = {
      SCHEMA_FILE    = "${path.module}/schema.json"
      PROJECT_REF    = supabase_project.ponderada.id
      SUPABASE_TOKEN = var.supabase_access_token
    }

    command = <<EOT
      python3 -c "
import json, os, urllib.request, urllib.error

with open(os.environ['SCHEMA_FILE']) as f:
    data = json.load(f)

sql = '\n'.join(stmt for m in data['migrations'] for stmt in m['sql'])

req = urllib.request.Request(
    f\"https://api.supabase.com/v1/projects/{os.environ['PROJECT_REF']}/database/query\",
    data=json.dumps({'query': sql}).encode(),
    headers={
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + os.environ['SUPABASE_TOKEN'],
        'User-Agent': 'terraform-supabase-provisioner/1.0'
    },
    method='POST'
)

try:
    with urllib.request.urlopen(req) as res:
        print('Schema aplicado com sucesso:', res.read().decode())
except urllib.error.HTTPError as e:
    print('Erro:', e.code, e.read().decode())
    raise
"
    EOT
  }
}