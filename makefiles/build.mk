##################################
### Building images
##################################

build.bite: build-bite ## Build docker image for Bite server
build-bite: executable
	@sh -c "./scripts/container_bite_build.sh"

