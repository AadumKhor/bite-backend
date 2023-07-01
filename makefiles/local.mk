#############################
### Local Setup for testing
############################
.ONESHELL:

executable:
	@chmod +x ${CURIDIR}/scripts/*

local.server: watch-server ## Start server for local testing
watch-server: 
	nodemon --watch cmd/ \
	--watch internal -e go --signal SIGKILL \
	--exec "go run ./cmd/main.go || exit 1"
