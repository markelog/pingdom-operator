.PHONY: install test build clean pack deploy ship revive golangci-lint sec vet

TAG ?= $(shell git rev-list HEAD --max-count=1 --abbrev-commit)
PRJ ?=
GO = GO111MODULE=on go
GO_FILES ?= ./...

install:
	@echo "[+] install"
	@go get $(GO_FILES)

test:
	@echo "[+] test"
	@go test $(GO_FILES)

build:
	@echo "[+] build"
	@go build -ldflags "-X main.version=$(TAG)" -o pingdom-operator .

clean:
	@echo "[+] clean"
	@rm ./pingdom-operator

pack:
	@echo "[+] pack"
	GOOS=linux make build
	docker build -t $(PRJ)/pingdom-operator:$(TAG) .

tag:
	@echo "[+] tag"
	@docker tag $(PRJ)/pingdom-operator:$(TAG) eu.gcr.io/$(PRJ)/pingdom-operator:$(TAG)

upload: pack tag
	@echo "[+] upload"
	@docker push eu.gcr.io/$(PRJ)/pingdom-operator:$(TAG)

deploy:
	@echo "[+] deploy"
	@kubectl apply -f ./examples

ship: test pack upload deploy clean

# ---

setup:
	@echo "[+] setup"
	@kubectl create -f deploy/service_account.yaml
	@kubectl create -f deploy/role.yaml
	@kubectl create -f deploy/role_binding.yaml
	@kubectl create -f deploy/operator.yaml


destroy:
	@echo "[+] destroy"
	@kubectl delete -f deploy/operator.yaml
	@kubectl delete -f deploy/role_binding.yaml
	@kubectl delete -f deploy/role.yaml
	@kubectl delete -f deploy/service_account.yaml

# Lint

lint: vet golangci-lint revive sec

scripts/bin/revive: scripts/go.mod
	@cd scripts; \
	$(GO) build -o ./bin/revive github.com/mgechev/revive

scripts/bin/golangci-lint: scripts/go.mod
	@cd scripts; \
	$(GO) build -o ./bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

scripts/bin/gosec: scripts/go.mod
	@cd scripts; \
	$(GO) build -o ./bin/gosec github.com/securego/gosec/cmd/gosec

revive: scripts/bin/revive
	@echo "lint via revive"
	scripts/bin/revive \
		-formatter stylish \
		-exclude ./vendor/... \
		-exclude pkg/apis/pingdom/v1alpha1/zz_generated.openapi.go \
		-exclude pkg/apis/pingdom/v1alpha1/zz_generated.deepcopy.go \
		$(GO_FILES)

golangci-lint: scripts/bin/golangci-lint
	@echo "lint via golangci-lint"
	@scripts/bin/golangci-lint run \
		--config ./scripts/configs/.golangci.yml \
		$(GO_FILES)

sec: scripts/bin/gosec
	@echo "lint via gosec"
	@scripts/bin/gosec -quiet \
		-exclude=G104,G107,G108,G201,G202,G204,G301,G304,G401,G402,G501 \
		-conf=./scripts/configs/gosec.json \
		$(GO_FILES)

vet:
	@echo "lint via go vet"
	@$(GO) vet $(GO_FILES)
