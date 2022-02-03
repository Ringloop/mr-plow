build:
	go build -o bin/mrplow cmd/main.go

test:
	go test ./...
clean:
	@echo "Cleaning the mrflow"
	@rm -fr bin/mrplow
