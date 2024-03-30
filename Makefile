docker-up:
	docker compose -f deployments/local/docker-compose.yml up -d
docker-ps:
	docker compose -f deployments/local/docker-compose.yml ps
docker-logs:
	docker compose -f deployments/local/docker-compose.yml logs 