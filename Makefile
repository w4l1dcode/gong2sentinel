
all: run

run:
	go run ./cmd/... -config=dev.yml

build:
	CGO_ENABLED=0 go build -o gong2sentinel ./cmd/...

test:
	go test -v ./...

clean:
	rm -r dist/ gong2sentinel || true

update:
	go get -u ./...
	go mod tidy
