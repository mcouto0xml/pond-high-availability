include consumer/.env
export

.PHONY: env-test
env-test:
	echo $(DATABASE_URL)

.PHONY: artefact-build
artefact-build:
	gcloud functions deploy PostTelemetry \
  --project=ponderada \
  --gen2 \
  --runtime go126 \
  --region us-central1 \
  --source=./consumer \
  --entry-point PostTelemetry \
  --trigger-http \
  --no-allow-unauthenticated \
  --set-env-vars DATABASE_URL=$(DATABASE_URL)