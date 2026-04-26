PROJ=fastlinks

.PHONY: startdb
startdb:
	docker compose -p $(PROJ) -f docker/docker-compose.yaml up -d --no-recreate

.PHONY: stopdb
stopdb:
	docker compose -p $(PROJ) -f docker/docker-compose.yaml down
