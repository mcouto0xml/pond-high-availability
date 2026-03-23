include consumer/.env
export


.PHONY: curl-producer
curl-producer:
	curl -X POST http://localhost:8080/telemetry \
		-H "Content-Type: application/json" \
		-d '{"iot_name": "sensor-01", "temperature": 23.5, "humidity": 60.2, "presence": true, "vibration": 0.04, "luminosity": 850.0, "tank_level": 75.3}'
	curl -X POST http://localhost:8080/telemetry \
		-H "Content-Type: application/json" \
		-d '{"iot_name": "sensor-01", "temperature": 24.5, "humidity": 68.3, "presence": true, "vibration": 0.06, "luminosity": 600.0, "tank_level": 72.3}'


.PHONY: curl-stress
curl-stress:
	for i in $$(seq 1 1000	); do \
		JSON=$$(python3 -c "import json; print(json.dumps({'iot_name': f'sensor-01', 'temperature': round(20 + $$i * 0.5, 2), 'humidity': round(60 + $$i * 0.8, 2), 'presence': $$i % 2 == 0, 'vibration': round(0.01 + $$i * 0.005, 3), 'luminosity': round(500 + $$i * 20.0, 2), 'tank_level': round(65 + $$i * 0.7, 2)}))"); \
		curl -s -o /dev/null -w "Request $$i: %{http_code}\n" -X POST http://localhost:8080/telemetry \
			-H "Content-Type: application/json" \
			-d "$$JSON"; \
	done