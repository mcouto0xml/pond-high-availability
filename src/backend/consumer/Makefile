.PHONY: artefact-build
artefact-build:
	gcloud functions deploy PostTelemetry \
  --project=ponderada \
  --gen2 \
  --runtime=go126 \
  --region=us-central1 \
  --source=. \
  --entry-point=PostTelemetry \
  --trigger-http \
  --no-allow-unauthenticated