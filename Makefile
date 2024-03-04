# Description: Makefile for the project

dev:
	@go run ./cmd/ui/main.go

gen:
	@go run ./cmd/generator/main.go

validate:
	@go run ./cmd/validator/main.go

.PHONY: dev gen validate