build:
	@mkdir -p bin
	go build -o ./bin/forgectl ./cmd/forgectl
	@echo "built: forgectl"
