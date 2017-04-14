REGISTRY ?= 
TAG ?= $(REGISTRY)/wotio/hydra
VERSION ?= latest
BUILD_DOCKERFILE ?= Dockerfile-wotio-build
DEPLOY_DOCKERFILE ?= Dockerfile-wotio-deploy

image:
	docker build --tag $(TAG):$(VERSION)-build --file $(BUILD_DOCKERFILE) .
	docker run -t -v `pwd`/bin:/deploy $(TAG):$(VERSION)-build cp /go/bin/hydra /deploy
	docker build --tag $(TAG):$(VERSION) --file $(DEPLOY_DOCKERFILE) .
