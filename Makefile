# =============================================================================

run-local:
	@go run app/services/bender/main.go

tidy:
	go mod tidy