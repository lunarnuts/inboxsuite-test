docker-compose-dev:
	go mod tidy && \
	go mod vendor && \
	docker-compose -f deploy/local/docker-compose.yml up --build
docker-compose:
	docker-compose -f deploy/local/docker-compose.yml up --build