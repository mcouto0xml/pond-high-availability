# Ponderada 2 — Backend de Telemetria Industrial com Alta Disponibilidade

Backend em **Go** para ingestão assíncrona de dados de sensores IoT industriais, utilizando **Google Cloud Tasks** como broker de mensagens, **Google Cloud Functions** como consumer e **PostgreSQL (Supabase)** como banco de dados relacional. Toda a infraestrutura em nuvem é provisionada via **Terraform**.

---

## Sumário

- [Contexto](#-contexto)
- [Arquitetura](#-arquitetura)
- [Estrutura do Repositório](#-estrutura-do-repositório)
- [Modelo de Dados](#-modelo-de-dados)
- [Decisões Técnicas](#-decisões-técnicas)
- [Pré-requisitos](#-pré-requisitos)
- [Como Rodar](#-como-rodar)
- [Endpoints da API](#-endpoints-da-api)
- [Teste de Carga](#-teste-de-carga)
- [Resultados e Análise](#-resultados-e-análise)

---

## Contexto

Uma empresa de monitoramento industrial precisa modernizar sua operação, coletando dados de sensores distribuídos em diferentes ambientes: temperatura, umidade, presença, vibração, luminosidade e nível de reservatórios. Com o crescimento do número de dispositivos conectados, surgem desafios de escalabilidade, confiabilidade e desempenho.

A solução adota uma **arquitetura desacoplada baseada em mensageria**, capaz de absorver picos de carga sem comprometer a estabilidade da aplicação — nenhum processamento pesado ocorre de forma síncrona no momento da requisição.

---

## Arquitetura

![Arquitetura utilizada](images/Cloud%20Architecture%20v1.png)

```
┌───────────────────────────────────────────────────────────────────┐
│                      Dispositivos IoT                             │
│             (sensores embarcados — múltiplos devices)             │
└────────────────────────────┬──────────────────────────────────────┘
                             │  HTTP POST /telemetry
                             ▼
┌───────────────────────────────────────────────────────────────────┐
│              Producer — Backend Go (HTTP :8080)                   │
│  Recebe o payload, serializa e despacha uma Task assíncrona       │
│  A resposta HTTP 202 é retornada imediatamente                    │
└────────────────────────────┬──────────────────────────────────────┘
                             │  CreateTask (HTTPS)
                             ▼
┌───────────────────────────────────────────────────────────────────┐
│              Google Cloud Tasks (Fila gerenciada)                 │
│  Retry automático · Rate limiting · Backoff exponencial           │
│  max_dispatches_per_second: 10 · max_concurrent_dispatches: 50    │
└────────────────────────────┬──────────────────────────────────────┘
                             │  HTTP POST (OIDC autenticado)
                             ▼
┌───────────────────────────────────────────────────────────────────┐
│              Consumer — Google Cloud Functions (Gen 2)            │
│  Runtime: Go 1.25 · Entry point: PostTelemetry                   │
│  Processa a mensagem e persiste no banco de dados                 │
└────────────────────────────┬──────────────────────────────────────┘
                             │  INSERT (go-pg)
                             ▼
┌───────────────────────────────────────────────────────────────────┐
│              PostgreSQL — Supabase                                 │
│  Tabelas: devices + telemetry (FK, índices, RLS)                  │
└───────────────────────────────────────────────────────────────────┘
```

### Fluxo detalhado

1. O dispositivo envia um `POST /telemetry` com payload JSON contendo os dados do sensor.
2. O **Producer** (Go HTTP server) valida e deserializa o payload, serializa a mensagem e cria uma **Task** no Google Cloud Tasks usando a API GRPC. A resposta `202 Accepted` é retornada de imediato — a chamada ao Cloud Tasks é disparada em uma goroutine separada para não bloquear o handler.
3. O **Cloud Tasks** gerencia a fila com retry automático (até 5 tentativas, backoff exponencial de 10s a 4000s) e despacha a tarefa via HTTP POST autenticado com OIDC para a URL do Consumer.
4. O **Consumer** (Cloud Function em Go) recebe a requisição, deserializa o payload e persiste o registro de telemetria no PostgreSQL — buscando o device por nome e inserindo na tabela `telemetry` via FK.
5. Em caso de falha no Consumer, o Cloud Tasks reencaminha a Task automaticamente, garantindo entrega sem perda de dados.

---

## Estrutura do Repositório

```
pond-high-availability/
├── src/
│   ├── backend/
│   │   ├── Makefile                        # Comandos utilitários (curl, stress test)
│   │   │
│   │   ├── producer/                       # API HTTP — recebe e enfileira telemetria
│   │   │   ├── cmd/run/main.go             # Entrypoint do servidor
│   │   │   ├── internal/
│   │   │   │   ├── config/
│   │   │   │   │   ├── cloudTasks.go       # Cliente do Google Cloud Tasks
│   │   │   │   │   └── queueInterface.go   # Interface da fila (desacoplamento)
│   │   │   │   ├── dto/
│   │   │   │   │   ├── telemetryNewDataRequest.go   # DTO de entrada
│   │   │   │   │   └── telemetryNewDataResponse.go  # DTO de saída
│   │   │   │   ├── endpoints/main.go       # Registro de rotas HTTP
│   │   │   │   ├── handlers/
│   │   │   │   │   ├── telemetryHandler.go # Handler POST /telemetry
│   │   │   │   │   └── healthHandler.go    # Handler GET /healthz
│   │   │   │   └── server/main.go          # Inicialização do servidor HTTP
│   │   │   ├── .env.example
│   │   │   ├── curl.example.txt
│   │   │   ├── go.mod
│   │   │   └── go.sum
│   │   │
│   │   └── consumer/                       # Cloud Function — processa e persiste
│   │       ├── postTelemetry.go            # Entrypoint da Cloud Function
│   │       ├── internal/
│   │       │   ├── db/postgreSql.go        # Repositório PostgreSQL (go-pg)
│   │       │   ├── dbContext/dbContext.go  # Interface do banco (desacoplamento)
│   │       │   ├── dto/consumerRequestDto.go
│   │       │   ├── models/
│   │       │   │   ├── devicesModel.go     # Model da tabela devices
│   │       │   │   └── telemetryModel.go   # Model da tabela telemetry
│   │       │   └── utils/loadEnvUtil.go
│   │       ├── .env.example
│   │       ├── curl.example.txt
│   │       ├── go.mod
│   │       └── go.sum
│   │
│   ├── infra/
│   │   ├── terraform/                      # IaC — provisiona toda a infra GCP
│   │   │   ├── main.tf                     # Orquestra os módulos
│   │   │   ├── variables.tf
│   │   │   ├── terraform.tfvars.example
│   │   │   └── modules/
│   │   │       ├── cloud_tasks/            # Fila + Service Account
│   │   │       ├── cloud_function/         # Consumer como Cloud Function Gen 2
│   │   │       ├── cloud_storage/          # Bucket para o source da Function
│   │   │       ├── postgresql/             # Supabase via Terraform + schema
│   │   │       └── iam/                    # Permissões e papéis GCP
│   │   └── gcp_oidc/
│   │       └── setup-wif.sh                # Script para Workload Identity Federation
│   │
│   └── test/
│       ├── main.go                         # Teste de carga em Go
│       └── go.mod
│
├── .gitignore
└── README.md
```

---

## Modelo de Dados

O banco utiliza duas tabelas com relacionamento 1:N — um dispositivo possui múltiplos registros de telemetria.

### `devices` — Cadastro dos dispositivos IoT

| Coluna        | Tipo        | Descrição                                  |
|---------------|-------------|--------------------------------------------|
| `id`          | BIGSERIAL   | Identificador único (PK, autoincrement)    |
| `name`        | TEXT        | Nome descritivo do dispositivo             |
| `description` | TEXT        | Localização ou descrição do device         |
| `created_at`  | TIMESTAMPTZ | Timestamp de cadastro (DEFAULT NOW())      |

### `telemetry` — Registros de leitura dos sensores

| Coluna        | Tipo        | Descrição                                        |
|---------------|-------------|--------------------------------------------------|
| `id`          | BIGSERIAL   | Identificador sequencial (PK)                    |
| `iot_id`      | BIGINT      | FK → `devices.id` (ON DELETE CASCADE)            |
| `temperature` | FLOAT8      | Temperatura em °C                                |
| `humidity`    | FLOAT8      | Umidade relativa do ar em %                      |
| `presence`    | BOOLEAN     | Detecção de presença (`true` = detectado)        |
| `vibration`   | FLOAT8      | Intensidade de vibração em m/s²                  |
| `luminosity`  | FLOAT8      | Luminosidade em lux (lx)                         |
| `tank_level`  | FLOAT8      | Nível do tanque em % (0.0 a 100.0)               |
| `created_at`  | TIMESTAMPTZ | Timestamp de inserção (DEFAULT NOW())            |

**Índices criados:**
- `idx_telemetry_created_at` — consultas por série temporal (DESC)
- `idx_telemetry_iot_id` — filtros por dispositivo

**Row Level Security (RLS)** está habilitado em ambas as tabelas para uso com Supabase.

O schema SQL completo está em `src/infra/terraform/modules/postgresql/schema.json`.

---

## Decisões Técnicas

### Go como linguagem principal

Go foi escolhido por seu excelente desempenho em cenários de alta concorrência, modelo de goroutines nativo e biblioteca padrão robusta para HTTP. A compilação estática facilita o deploy em Cloud Functions sem dependências externas.

Um detalhe relevante na implementação: a chamada ao Cloud Tasks no handler de telemetria é disparada em uma **goroutine separada** (`go func() { ... }()`), o que fez uma diferença significativa na latência percebida pelo cliente — o `202 Accepted` retorna imediatamente sem aguardar a confirmação de enfileiramento.

### Google Cloud Tasks como broker de mensagens

Em vez de manter um broker autogerenciado (como RabbitMQ), optou-se pelo **Cloud Tasks** — serviço gerenciado do GCP. Isso traz:

- **Retry automático** com backoff exponencial configurável (10s a 4000s, até 5 tentativas)
- **Rate limiting** nativo (10 dispatches/s, 50 concorrentes) para proteger o consumer
- **Zero operação** — sem necessidade de provisionar, monitorar ou escalar o broker
- **Autenticação OIDC** entre o Cloud Tasks e a Cloud Function, sem exposição de endpoints públicos desprotegidos

### Google Cloud Functions (Gen 2) como Consumer

O consumer é uma Cloud Function em Go, o que oferece:

- **Scale-to-zero** — sem custos quando não há mensagens
- **Scale-out automático** até 50 instâncias
- Inicialização com `sync.Once` garante que a conexão com o banco seja criada apenas uma vez por instância, sendo reutilizada em todas as invocações subsequentes

### PostgreSQL via Supabase

Supabase fornece PostgreSQL gerenciado com RLS, API REST e dashboard integrados. A conexão é feita via URL com TLS. O ORM utilizado é o **go-pg/v10**, que suporta o padrão de model-based queries do Go.

### Terraform para Infraestrutura como Código

Toda a infra GCP — Cloud Tasks, Cloud Functions, Cloud Storage, IAM e Supabase — é provisionada e versionada via Terraform, organizado em módulos reutilizáveis. Isso garante reproducibilidade total do ambiente e facilita auditorias.

### Teste de carga em Go (em vez de k6)

Apesar de o enunciado recomendar o k6, o teste de carga foi implementado em **Go nativo**. A escolha foi deliberada: aproveitar o mesmo ecossistema da aplicação, eliminar dependências externas e explorar as goroutines do Go para simular concorrência real. O resultado é um binário único, configurável por flags, sem necessidade de instalar ferramentas adicionais.

---

## Pré-requisitos

- [Go 1.21+](https://go.dev/dl/)
- [Terraform 1.5+](https://developer.hashicorp.com/terraform/install)
- Conta no [Google Cloud Platform](https://cloud.google.com/) com projeto criado
- Conta no [Supabase](https://supabase.com/) com organização criada
- `gcloud` CLI autenticado (`gcloud auth application-default login`)

---

## Como Rodar

### 1. Clonar o repositório

```bash
git clone https://github.com/mcouto0xml/pond-high-availability.git
cd pond-high-availability
```

### 2. Provisionar a infraestrutura com Terraform

```bash
cd src/infra/terraform

# Copiar e preencher as variáveis
cp terraform.tfvars.example terraform.tfvars
# Edite terraform.tfvars com seus valores reais

# Inicializar e aplicar
terraform init
terraform plan
terraform apply
```

Ao final do `apply`, anote os outputs: URL da Cloud Function (será o `WORKER_URL` do producer) e o email da Service Account do Cloud Tasks.

As variáveis necessárias em `terraform.tfvars`:

| Variável                | Descrição                                              |
|-------------------------|--------------------------------------------------------|
| `project_id`            | ID do projeto GCP                                      |
| `region`                | Região GCP (ex: `us-central1`)                         |
| `tasks_name`            | Nome da fila no Cloud Tasks                            |
| `admin_email`           | Email do admin com permissão no Cloud Tasks            |
| `supabase_access_token` | Token de acesso do Supabase                            |
| `supabase_org_id`       | ID da organização no Supabase                          |
| `database_region`       | Região do banco Supabase (ex: `us-east-1`)             |
| `database_password`     | Senha do banco PostgreSQL                              |
| `function_name`         | Nome da Cloud Function consumer                        |
| `function_source_dir`   | Caminho do código do consumer (relativo ao terraform/) |
| `environment_variables` | Env vars injetadas na Cloud Function (DATABASE_URL)    |

### 3. Configurar o Schema do banco

Após o provisionamento do Supabase, aplique o schema SQL disponível em `src/infra/terraform/modules/postgresql/schema.json`. Copie os SQLs do campo `migrations[0].sql` e execute no SQL Editor do Supabase Dashboard, ou via `psql`:

```bash
psql "postgresql://<user>:<password>@<host>:5432/postgres" -f schema.sql
```

Registre manualmente os devices na tabela `devices` (por nome, pois o consumer faz lookup por `iot_name`):

```sql
INSERT INTO public.devices (name, description) VALUES
  ('sensor-01', 'Sensor de temperatura do galpão A'),
  ('sensor-02', 'Sensor de umidade do reservatório B'),
  ('sensor-03', 'Sensor de presença da entrada principal');
```

### 4. Configurar e subir o Producer

```bash
cd src/backend/producer

cp .env.example .env
# Edite .env com os valores do seu projeto GCP
```

Conteúdo do `.env`:

```env
PROJECT_ID="seu-projeto-gcp"
QUEUE_LOCATION="us-central1"
QUEUE_ID="nome-da-sua-fila"
WORKER_URL="https://url-da-sua-cloud-function"
SERVICE_ACCOUNT_EMAIL="cloud-task-editor@seu-projeto.iam.gserviceaccount.com"
```

Rodando localmente:

```bash
go run cmd/run/main.go
```

O servidor sobe na porta `:8080`.

### 5. Testar a API manualmente

```bash
curl -X POST http://localhost:8080/telemetry \
  -H "Content-Type: application/json" \
  -d '{
    "iot_name": "sensor-01",
    "temperature": 23.5,
    "humidity": 60.2,
    "presence": true,
    "vibration": 0.04,
    "luminosity": 850.0,
    "tank_level": 75.3
  }'
```

Resposta esperada (`202 Accepted`):
```json
{ "message": "Mensagem adicionada a fila!" }
```

### 6. Comandos via Makefile

```bash
cd src/backend

# Envia duas requisições de exemplo ao producer
make curl-producer

# Executa um stress test via shell (1000 requisições sequenciais)
make curl-stress
```

---

## Endpoints da API

### `GET /healthz`

Verifica se o servidor está operacional.

**Resposta `200 OK`:**
```json
{ "status": "ok" }
```

---

### `POST /telemetry`

Recebe um pacote de telemetria de um dispositivo embarcado e o enfileira no Cloud Tasks para processamento assíncrono.

**Request Body:**

```json
{
  "iot_name": "sensor-01",
  "temperature": 23.5,
  "humidity": 60.2,
  "presence": true,
  "vibration": 0.04,
  "luminosity": 850.0,
  "tank_level": 75.3
}
```

| Campo         | Tipo    | Descrição                                              |
|---------------|---------|--------------------------------------------------------|
| `iot_name`    | string  | Nome do dispositivo (deve existir na tabela `devices`) |
| `temperature` | float64 | Temperatura em °C                                      |
| `humidity`    | float64 | Umidade relativa em %                                  |
| `presence`    | bool    | Presença detectada (`true`/`false`)                    |
| `vibration`   | float64 | Intensidade de vibração em m/s²                        |
| `luminosity`  | float64 | Luminosidade em lux                                    |
| `tank_level`  | float64 | Nível do tanque em % (0.0 a 100.0)                     |

**Responses:**

| Status | Descrição                                        |
|--------|--------------------------------------------------|
| `202`  | Payload aceito e Task criada na fila com sucesso |
| `400`  | Payload inválido ou erro de deserialização       |
| `405`  | Método não permitido (somente POST)              |

---

## Teste de Carga

O teste de carga foi implementado em **Go nativo** (`src/test/main.go`), utilizando goroutines e canais para simular múltiplos dispositivos enviando requisições simultâneas. Apesar de o enunciado recomendar o k6, optou-se pelo Go para manter o mesmo ecossistema da aplicação e eliminar dependências externas.

### Como executar

```bash
cd src/test
go run main.go [flags]
```

### Flags disponíveis

| Flag        | Padrão                              | Descrição                          |
|-------------|-------------------------------------|------------------------------------|
| `-url`      | `http://localhost:8080/telemetry`   | URL alvo                           |
| `-c`        | `10`                                | Goroutines concorrentes            |
| `-n`        | `10000`                             | Total de requisições               |
| `-timeout`  | `10s`                               | Timeout por requisição HTTP        |

**Exemplos:**

```bash
# Teste padrão: 10 goroutines, 10.000 requisições
go run main.go

# Teste de stress: 100 goroutines, 50.000 requisições
go run main.go -c 100 -n 50000

# Apontando para produção
go run main.go -url https://minha-api.cloud.com/telemetry -c 50 -n 20000
```

### O que o teste faz

O test runner:
1. Enfileira `N` jobs em um channel bufferizado.
2. Sobe `C` goroutines trabalhadoras que consomem jobs do channel, cada uma gerando um payload aleatório e válido (sensores `sensor-01`, `sensor-02` e `sensor-03` com valores realistas).
3. Coleta resultados de forma thread-safe com mutex + atomic counters.
4. Exibe progresso a cada 2 segundos via goroutine de monitoramento.
5. Imprime relatório final ao término.

**Saída do relatório:**

```
╔══════════════════════════════════════════╗
║         LOAD TEST RESULTS                ║
╚══════════════════════════════════════════╝

📊 General
   Total Requests  : 1000
   Elapsed Time    : 132ms
   Req/s (RPS)     : 7568.08

✅ Results
   Success (2xx)   : 1000 (100%)
   Failed          : 0 (0%)

⏱️  Latency
   Min             : 309µs
   Avg             : 1.311µs
   Max             : 9.431ms

📋 Status Codes
   HTTP 202         : 1000
```

---

## 📊 Resultados e Análise

> Teste executado localmente com o Producer rodando em `localhost:8080`, apontando para a infraestrutura real no GCP.

### Configuração do teste

| Parâmetro     | Valor    |
|---------------|----------|
| Goroutines    | 10       |
| Total de reqs | 1000   |
| Timeout       | 10s      |
| Sensores      | 3 devices |

### Métricas obtidas

| Métrica                 | Valor         |
|-------------------------|---------------|
| Total de requisições    | 1000        |
| Sucesso (2xx)           | ~100%        |
| Taxa de erro            | ~0%         |
| RPS (throughput)        | ~7568.08 req/s    |
| Latência mínima         | ~309ms          |
| Latência média          | ~1.311ms         |
| Latência máxima         | ~9.431ms        |


### Análise

**Throughput:** A API conseguiu absorver o pico de carga com baixíssima taxa de erro. O desacoplamento via Cloud Tasks é o principal responsável: o handler retorna imediatamente após enfileirar, sem aguardar processamento no banco.

**Latência:** A latência média baixa (~1.311ms) reflete que o gargalo não está no Producer, mas potencialmente no Cloud Tasks → Cloud Function → Supabase, que opera de forma totalmente assíncrona e invisível ao cliente.

**Goroutine no handler:** A chamada `go func() { cloudTasks.CreateTask(...) }()` no handler eliminou o tempo de espera pela resposta da API do GCP no caminho crítico da requisição, reduzindo significativamente a latência percebida pelo dispositivo.

### Gargalos identificados

**Cloud Tasks rate limit:** A fila está configurada com `max_dispatches_per_second: 10`, o que protege o consumer de sobrecarga mas limita o throughput de processamento. Em produção com alto volume, aumentar esse limite ou usar múltiplas filas pode ser necessário.

**Consumer stateless vs. conexões com o banco:** A Cloud Function usa `sync.Once` para reutilizar a conexão com o PostgreSQL por instância. Em scale-out com muitas instâncias simultâneas, pode-se atingir o limite de conexões do Supabase. O uso de um connection pooler (ex: Supabase Pooler com PgBouncer) mitiga esse risco.

**Goroutine fire-and-forget no Producer:** A chamada ao Cloud Tasks é feita em goroutine sem tratamento de erro retornado ao cliente. Se a Task falhar ao ser criada (ex: quota esgotada, erro de rede), o cliente recebe `202` mas a mensagem é perdida. Uma melhoria seria implementar um mecanismo de fallback ou log estruturado de erros.

### Melhorias possíveis

- Aumentar `max_dispatches_per_second` da fila conforme a capacidade do Consumer.
- Adicionar Dead Letter Queue (DLQ) para Tasks que excedem o número máximo de tentativas.
- Implementar observabilidade com Cloud Monitoring + alertas de fila crescente.
- Adicionar índice composto `(iot_id, created_at)` para consultas analíticas por dispositivo em janelas de tempo.
- Configurar Supabase Pooler para gerenciar conexões em cenários de alta concorrência.
