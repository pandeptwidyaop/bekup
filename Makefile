docker-up:
	docker compose -p bekup-development  -f deployments/local/docker-compose.yml up -d
docker-down:
	docker-compose -p bekup-development -f deployments/local/docker-compose.yml down 
docker-ps:
	docker compose -p bekup-development -f deployments/local/docker-compose.yml ps
docker-logs:
	docker compose -p bekup-development -f deployments/local/docker-compose.yml logs
docker-rebuild:
	docker compose -p bekup-development -f deployments/local/docker-compose.yml up -d --build
docker-exec:
	docker compose -p bekup-development -f deployments/local/docker-compose.yml exec app bash
build: