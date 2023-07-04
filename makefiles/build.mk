##################################
### Building images / entire setup
##################################

build.bite: build-bite ## Build docker image for Bite server
build-bite: executable
	@sh -c "./scripts/container_bite_build.sh"

build: build-all ## Build app & DB both without cache
build-all: remove-docker-compose-images docker-compose-up

build.all.cache: docker-compose-up ## Build app & DB with cache
docker-compose-up:
	@sh -c "sudo docker compose up --force-recreate"

stop: docker-compose-down ## Stop app & DB both
docker-compose-down: 
	@sh -c "sudo docker compose down"

image.clear: remove-docker-compose-images ## Remove image to re-build from scratch
remove-docker-compose-images: 
	@sh -c "sudo docker image rm bitespeed-backend-task-app --force || true"

build.db: build-db ## Starting just the DB in detached state
build-db: 
	@sh -c "sudo docker compose up db -d"
