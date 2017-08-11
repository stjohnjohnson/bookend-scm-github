build:
	go build -a -ldflags "-X main.VERSION=${VERSION}"

test:
	go test -cover ./...
