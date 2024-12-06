# ==============================================================================
# Define dependencies

KIND            := kindest/node:v1.29.1@sha256:a0cc28af37cf39b019e2b448c54d1a3f789de32536cb5a5db61a49623e527144
KIND_CLUSTER    := dev
NAMESPACE       := bender-system
APP             := bender
BASE_IMAGE_NAME := ghcr.io/zmoog
SERVICE_NAME    := bender
VERSION:= 0.0.1-$(shell git rev-parse --short HEAD)
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

# =============================================================================
# Building containers

service:
	docker build \
		-f zarf/docker/dockerfile.service \
		-t ${SERVICE_IMAGE} \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
		.

# =============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image=${KIND} \
		--name=${KIND_CLUSTER} \
		--config=zarf/k8s/dev/kind-config.yaml


dev-down:
	kind delete cluster --name ${KIND_CLUSTER}

dev-load:
	kind load docker-image ${SERVICE_IMAGE} --name ${KIND_CLUSTER}

dev-apply:
	kustomize build zarf/k8s/dev/bender | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --for=condition=Ready --timeout=60s

# =============================================================================

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-update: service dev-load dev-restart

dev-update-apply: service dev-load dev-apply

# =============================================================================

dev-logs:
	@kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true --tail=100 -f | go run app/tooling/logfmt/main.go

dev-describe-deployment:
	@kubectl describe deployment --namespace=$(NAMESPACE) $(APP)

dev-describe-bender:
	@kubectl describe pods --namespace=$(NAMESPACE) -l app=$(APP)

# =============================================================================

run-local:
	go run app/services/bender-bot/main.go | go run app/tooling/logfmt/main.go

run-local-help:
	go run app/services/bender/main.go --help

tidy:
	go mod tidy

lint:
	golangci-lint run

test:
	go test -v ./...
