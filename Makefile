.PHONY: install test build clean pack deploy ship

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
PRJ?=

export TAG

install:
	@echo "[+] install"
	@go get ./...

test:
	@echo "[+] test"
	@go test ./...

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
