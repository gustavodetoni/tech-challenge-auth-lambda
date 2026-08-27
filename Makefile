build:
	mkdir -p bin
	go build -o ./bin/auth ./cmd

tests:
	mkdir -p /tmp/go-tmp /tmp/go-build-cache-auth-lambda /tmp/go-mod-cache-auth-lambda
	env GOCACHE=/tmp/go-build-cache-auth-lambda GOTMPDIR=/tmp/go-tmp GOMODCACHE=/tmp/go-mod-cache-auth-lambda \
		go test ./... -count=1

package:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/bootstrap ./cmd
	cd dist && zip function.zip bootstrap

fmt:
	goimports -w .
