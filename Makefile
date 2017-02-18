
.PHONY: all
all:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/wp

.PHONY: install
install:
	sudo cp -rf bin/wp /usr/local/bin/wp

.PHONY: test
test:
	go vet .
	go test -race -test.v .


