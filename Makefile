build:
	@docker compose --env-file .env -f docker-compose.yml build

run:
	@docker compose --env-file .env -f docker-compose.yml up -d

down:
	@docker compose --env-file .env -f docker-compose.yml down

test:
	cd src && go test ./tests/*
	cd ../
