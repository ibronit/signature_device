run:
	@go run ./cmd/signature_service

tidy:
	@go mod tidy

test:
	@go test ./...