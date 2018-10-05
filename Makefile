all: push

APP?=ui
PROJECT?=github.com/k8s-community/${APP}
REGISTRY?=gcr.io/ws-aug-16
CA_DIR?=certs

# Use the 0.0.0 tag for testing, it shouldn't clobber any release builds
RELEASE?=0.5.11
GOOS?=linux
GOARCH?=amd64

K8SAPP_LOCAL_PORT?=8080
K8SAPP_LOCAL_HOST?=0.0.0.0
K8SAPP_LOG_LEVEL?=debug

DB_CONNECTION_STRING?=

NAMESPACE?=k8s-community
INFRASTRUCTURE?=stable
KUBE_CONTEXT?=${INFRASTRUCTURE}
VALUES?=values-${INFRASTRUCTURE}

CONTAINER_NAME?=${NAMESPACE}-${APP}
CONTAINER_IMAGE?=${REGISTRY}/${CONTAINER_NAME}

REPO_INFO=$(shell git config --get remote.origin.url)

ifndef COMMIT
	COMMIT := git-$(shell git rev-parse --short HEAD)
endif

BUILDTAGS=

.PHONY: all
all: build

.PHONY: build
build: clean certs
	@echo "+ $@"
	@CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -installsuffix cgo \
		-ldflags "-s -w -X ${PROJECT}/version.RELEASE=${RELEASE} -X ${PROJECT}/version.COMMIT=${COMMIT} -X ${PROJECT}/version.REPO=${REPO_INFO}" \
		-o bin/${GOOS}-${GOARCH}/${APP} ${PROJECT}/cmd
	docker build --pull -t $(CONTAINER_IMAGE):$(RELEASE) .

.PHONY: certs
certs:
ifeq ("$(wildcard $(CA_DIR)/ca-certificates.crt)","")
	@echo "+ $@"
	@docker run --name ${CONTAINER_NAME}-certs -d alpine:edge sh -c "apk --update upgrade && apk add ca-certificates && update-ca-certificates"
	@docker wait ${CONTAINER_NAME}-certs
	@mkdir -p ${CA_DIR}
	@docker cp ${CONTAINER_NAME}-certs:/etc/ssl/certs/ca-certificates.crt ${CA_DIR}
	@docker rm -f ${CONTAINER_NAME}-certs
endif

.PHONY: push
push: build
	@echo "+ $@"
	@docker push $(CONTAINER_IMAGE):$(RELEASE)

.PHONY: run
run: build
	@echo "+ $@"
	docker run --name ${CONTAINER_NAME} -p ${K8SAPP_LOCAL_PORT}:${K8SAPP_LOCAL_PORT} \
		-e "SERVICE_HOST=${K8SAPP_LOCAL_HOST}" \
		-e "SERVICE_PORT=${K8SAPP_LOCAL_PORT}" \
		-d $(CONTAINER_IMAGE):$(RELEASE)
	sleep 1
	docker logs ${CONTAINER_NAME}

HAS_RUNNED := $(shell docker ps | grep ${CONTAINER_NAME})
HAS_EXITED := $(shell docker ps -a | grep ${CONTAINER_NAME})

.PHONY: logs
logs:
	docker logs ${CONTAINER_NAME}

.PHONY: stop
stop:
ifdef HAS_RUNNED
	@echo "+ $@"
	@docker stop ${CONTAINER_NAME}
endif

.PHONY: start
start: stop
	@echo "+ $@"
	@docker start ${CONTAINER_NAME}

.PHONY: rm
rm:
ifdef HAS_EXITED
	@echo "+ $@"
	@docker rm ${CONTAINER_NAME}
endif

.PHONY: deploy
deploy: push
	helm upgrade ${CONTAINER_NAME} -f charts/${VALUES}.yaml charts \
		--kube-context ${KUBE_CONTEXT} --namespace ${NAMESPACE} --version=${RELEASE} -i --wait \
		--set image.registry=${REGISTRY} --set image.name=${CONTAINER_NAME} --set image.tag=${RELEASE}

.PHONY: fmt
fmt:
	@echo "+ $@"


.PHONY: lint
lint: bootstrap
	@echo "+ $@"

.PHONY: vet
vet:
	@echo "+ $@"
	@go vet ./...

.PHONY: test
test: clean fmt lint vet
	@echo "+ $@"
	@go test -v -race -tags "$(BUILDTAGS) cgo" ./...

.PHONY: clean
clean: stop rm
	@rm -f bin/${GOOS}-${GOARCH}/${APP}

HAS_DEP := $(shell command -v dep;)
HAS_LINT := $(shell command -v golint;)

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif
ifndef HAS_LINT
	go get -u github.com/golang/lint/golint
endif
