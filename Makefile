#.PHONY: help

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

## DOCKER
CONTAINER_NAME=appto-catalog

build: ## build and up docker containers
	@docker-compose up --build -d

start: ## run cart server
	@docker-compose up -d

make sh:
	@docker exec -it $(CONTAINER_NAME) sh

##TEST
test: ## run tests
	@docker exec -it $(CONTAINER_NAME) go test ./... -v
	#@go test ./...
	#go test ./listing -run TestProductLister_List -v
	#go test ./pricing -run TestPricingCalculater_Calculate -v

coverage:
	mkdir -p var/test_results
	@go test -coverprofile=var/test_results/coverage.out ./...
	@go tool cover -func=var/test_results/coverage.out

