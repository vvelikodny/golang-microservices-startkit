all: build-services

build-services:
	make -C query-client-service
	make -C storage-service

run-env: build-services
	@docker-compose up -d

stop:
	@docker-compose stop
	@docker-compose rm

deploy-local: run-env

.PHONY: deploy-local
