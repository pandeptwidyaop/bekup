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
	go build -ldflags="-s -w" -o bin/bekup cmd/main.go
docker-build:
	docker build -t  pandeptwidyaop/bekup -f deployments/distribute/Dockerfile .
docker-release:
	docker push pandeptwidyaop/bekup