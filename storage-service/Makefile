PROJECT=storage-service
HUB=registry.gitlab.com
IMAGE=vvelikodny/golang-microservices-test/storage-service

all: deploy-local

GOOS=linux
GOARCH=amd64
GO111MODULE=on
GOFLAGS=-mod=vendor

clean:
	[ ! -f ./$(PROJECT) ] || rm ./$(PROJECT) || echo "-"

$(PROJECT): clean
	echo Building for $(GOOS)/$(GOARCH)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=$(GO111MODULE) GOFLAGS=$(GOFLAGS) go build -a -o ./$(PROJECT)

docker-build: clean
	@docker build -t $(HUB)/$(IMAGE) .

deploy-local: docker-build

.PHONY: deploy-local
