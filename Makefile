GOLDFLAGS="-X main.branch=main.commit $(COMMIT)"

default: complete-build

get:
	@go get -d ./...

complete-build: get
	@go build -ldflags=$(GOLDFLAGS) -a -o bin/godutch ./cmd/godutch

build:
	@go build -o bin/godutch ./cmd/godutch
