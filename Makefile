
GOPACKAGES := $(shell go list ./... | grep -v /vendor/)

.PHONY: all
all:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/wp

.PHONY: test
test:
	go vet ${GOPACKAGES}
	go test -race -test.v ${GOPACKAGES}


