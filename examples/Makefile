gen:
	rm -rf *_gen.go
	go generate
	go vet ./...

test:
	go test -v ./...
