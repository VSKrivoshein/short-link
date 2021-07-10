default:
	go build -v ./cmd/apiserver && ./apiserver

test:
	go test -v -race -timeout 10s ./...