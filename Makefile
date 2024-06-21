setup_and_run_local:
	docker build --tag muzz-api .
	docker tag muzz-api:latest muzz-api:local
	docker-compose up -d