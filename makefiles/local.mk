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
