.PHONY: default test test-cover dev

# for dev
dev:
	fresh

# for test
test:
	go test -race -cover ./...

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

build-web:
	cd web \
		&& npm run build \
		&& cd .. \
		&& packr -z

bench:
	go test -bench=. ./...

build-test:
	go run tool/main.go -max 200

build:
	CGO_ENABLED=0 go run tool/main.go && go build -tags netgo -o location
