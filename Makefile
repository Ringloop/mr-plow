build:
	go build -o bin/mrplow cmd/main.go

test:
	go test -v  ./...
clean:
	@echo "Cleaning the mrflow"
	@rm -fr bin/mrplow
