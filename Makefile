.PHONY: clean coverage help test report pkgsite vuln

help: ## list available targets
	@# Shamelessly stolen from Gomega's Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

clean: ## cleans up build and testing artefacts
	rm -f coverage.html coverage.out

coverage: ## gathers coverage and updates README badge
	@scripts/cov.sh

pkgsite: ## serves Go documentation on port 6060
	@echo "navigate to: http://localhost:6060/github.com/thediveo/gitloaderfs"
	@scripts/pkgsite.sh

test: ## runs all tests
	go test -v -p 1 -count 1 ./...

report: ## runs goreportcard
	@scripts/goreportcard.sh

vuln: ## runs Go vulnerability checks
	@scripts/vuln.sh
