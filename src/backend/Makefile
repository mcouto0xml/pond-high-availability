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

.PHONY: curl-producer
.PHONY: curl-producer
curl-producer:
	curl -X POST http://localhost:8080/telemetry \
		-H "Content-Type: application/json" \
		-d '{"iot_name": "sensor-01", "temperature": 23.5, "humidity": 60.2, "presence": true, "vibration": 0.04, "luminosity": 850.0, "tank_level": 75.3}'

