# Aqui utilizei o Cloud para habilitar e me ensinar como se utiliza o OpenID Connect
# =============================================================================
# Setup: Workload Identity Federation (WIF) para GitHub Actions + GCP
# =============================================================================
#

set -euo pipefail

# ─── CONFIGURAÇÕES ─────────────────────────────────────────────────────────────
PROJECT_ID="ponderada"               # ID do projeto GCP
GITHUB_ORG="mcouto0xml"           # Org ou usuário do GitHub
GITHUB_REPO="pond-high-availability"         # Nome do repositório (sem org)

SA_NAME="github-actions-sa"              # Nome da Service Account
POOL_ID="github-pool"                    # ID do Workload Identity Pool
PROVIDER_ID="github-provider"            # ID do Provider
REGION="us-central1"              # Região do Artifact Registry
# ───────────────────────────────────────────────────────────────────────────────

PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format="value(projectNumber)")

echo "🔧 Configurando projeto: $PROJECT_ID (number: $PROJECT_NUMBER)"

# 1. Ativar APIs necessárias
echo "📡 Ativando APIs..."
gcloud services enable iamcredentials.googleapis.com \
  iam.googleapis.com \
  artifactregistry.googleapis.com \
  --project="$PROJECT_ID"

# 2. Criar a Service Account
echo "👤 Criando Service Account..."
gcloud iam service-accounts create "$SA_NAME" \
  --display-name="GitHub Actions SA" \
  --description="Usada pelo GitHub Actions via WIF" \
  --project="$PROJECT_ID"

SA_EMAIL="${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

# 3. Conceder permissão de Artifact Registry Writer para a SA
echo "🔑 Concedendo permissões no Artifact Registry..."
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:${SA_EMAIL}" \
  --role="roles/artifactregistry.writer"

# 4. Criar o Workload Identity Pool
echo "🏊 Criando Workload Identity Pool..."
gcloud iam workload-identity-pools create "$POOL_ID" \
  --location="global" \
  --display-name="GitHub Actions Pool" \
  --description="Pool para autenticação do GitHub Actions" \
  --project="$PROJECT_ID"

# 5. Criar o Workload Identity Provider (OIDC do GitHub)
echo "🔗 Criando Workload Identity Provider..."
gcloud iam workload-identity-pools providers create-oidc "$PROVIDER_ID" \
  --location="global" \
  --workload-identity-pool="$POOL_ID" \
  --display-name="GitHub Provider" \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository,attribute.repository_owner=assertion.repository_owner" \
  --attribute-condition="assertion.repository_owner == '${GITHUB_ORG}'" \
  --project="$PROJECT_ID"

# 6. Permitir que o repositório específico impersone a SA
echo "🤝 Vinculando repositório à Service Account..."
gcloud iam service-accounts add-iam-policy-binding "${SA_EMAIL}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/${POOL_ID}/attribute.repository/${GITHUB_ORG}/${GITHUB_REPO}" \
  --project="$PROJECT_ID"

# 7. Exibir os valores para usar no workflow
WORKLOAD_IDENTITY_PROVIDER="projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/${POOL_ID}/providers/${PROVIDER_ID}"

echo ""
echo "✅ WIF configurado com sucesso!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Adicione estes secrets no GitHub Actions:"
echo "   Settings → Secrets and variables → Actions"
echo ""
echo "   GCP_PROJECT_ID            = ${PROJECT_ID}"
echo "   GCP_WORKLOAD_IDENTITY_PROVIDER = ${WORKLOAD_IDENTITY_PROVIDER}"
echo "   GCP_SERVICE_ACCOUNT       = ${SA_EMAIL}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"