LDFLAGS=-ldflags "-X main.VERSION=${VERSION}"

build:
	go build -a ${LDFLAGS}

test:
	go test -cover ./...

install:
	go install ${LDFLAGS}
