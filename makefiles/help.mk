########################
### Helper
########################

.PHONY: help
help: ## Display this help
	@./makefiles/banner.sh
	@IFS=$$'\n'; for line in $$(grep -h -E '^[a-zA-Z_#-]+:?.*?## .*$$' $(MAKEFILE_LIST)); do if [ "$${line:0:2}" = "##" ]; then \
 	echo $$line | awk 'BEGIN {FS = "## "}; {printf "\n\033[33m%s\033[0m\n", $$2}'; else \
 	echo $$line | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m   %s\n", $$1, $$2}'; fi; \
 	done; unset IFS;