beforedo: test
	go build .

test:
	go test ./...

setup:
	go get ./...
