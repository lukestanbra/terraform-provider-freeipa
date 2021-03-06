TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=hashicorp.com
NAMESPACE=lukestanbra
NAME=freeipa
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=darwin_amd64

.EXPORT_ALL_VARIABLES:

export FREEIPA_HOST=ipa.example.test
export FREEIPA_USERNAME=admin
export FREEIPA_PASSWORD=password

default: install

container:
	if ! grep ipa.example.test /etc/hosts >/dev/null; then echo "127.0.0.1 ipa.example.test" >> /etc/hosts; fi
	docker compose up -d

certificate:
	docker cp freeipa-server:/etc/ipa/ca.crt ./ca.crt
	sudo security add-trusted-cert -d -r trustRoot -p ssl -k /Library/Keychains/System.keychain ca.crt

build:
	go mod tidy
	go get
	go mod vendor
	go build -o ${BINARY}

release:
	go generate
	goreleaser release --rm-dist --snapshot --skip-publish  --skip-sign

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

example:
	cd examples && rm -rf .terraform* && terraform init && terraform apply