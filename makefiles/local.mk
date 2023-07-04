#############################
### Local Setup for testing
############################
.ONESHELL:

executable:
	@chmod +x ${CURDIR}/scripts/*
	@chmod +x ${CURDIR}/makefiles/*

local.bite: watch-bite ## Start bite server for local testing
watch-bite: 
	nodemon --watch cmd/ \
	--watch internal -e go --signal SIGKILL \
	--exec "clear && go run ./cmd/main.go || exit 1"

run.db: run-db ## Running the DB container in detached state
run-db:
	@sh -c "sudo docker compose up db -d"

setup.db: run-db setup-db ## Setup DB from scratch (this will remove data from existing tables)
setup-db: 
	source env/db.envrc && sudo docker compose exec -T db psql -U $$DB_USER -d $$DB_NAME < $$SETUP_SCRIPT_PATH

seed.db: run-db seed-db ## Add some data to the DB for testing
seed-db: 
	source env/db.envrc && sudo docker compose exec -T db psql -U $$DB_USER -d $$DB_NAME < $$SEED_SCRIPT_PATH

run.all: run-all ## Build & run the app & DB both without using image cache
run-all: remove-docker-compose-images docker-compose-up

run.all.cache: docker-compose-up ## Build & run app & DB with cache
docker-compose-up:
	@sh -c "sudo docker compose up --force-recreate"

stop: docker-compose-down ## Stop app & DB both
docker-compose-down: 
	@sh -c "sudo docker compose down"

image.clear: remove-docker-compose-images ## Remove image to re-build from scratch
remove-docker-compose-images: 
	@sh -c "sudo docker image rm bitespeed-backend-task-app --force || true"