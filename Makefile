GOLDFLAGS="-X main.branch=main.commit $(COMMIT)"

default: complete-build

tmpdir:
	@mkdir -p /tmp/godutch

get:
	@go get -d ./...

complete-build: tmpdir get
	@go build -ldflags=$(GOLDFLAGS) -a -o bin/godutch ./cmd/godutch

build:
	@go build -o bin/godutch ./cmd/godutch
