HOSTNAME=github.com
NAMESPACE=ibm-hyper-protect
NAME=terraform-provider-ibm-cloud-hyper-protect-virtual-server-for-ibm-cloud-vpc
BINARY=terraform-provider-hpcr
VERSION=0.0.1
OS_ARCH=linux_amd64 #Need to change to your current architecture

default: install

doc:
	go generate

build:
	go build -o ${BINARY}

release:
	GOOS=linux   GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux   GOARCH=s390x go build -o ./bin/${BINARY}_${VERSION}_linux_s390x
	GOOS=darwin  GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=darwin  GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_darwin_arm64
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	cp -a ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test ./...

debug: build
	./${BINARY} -debug
