CURRENT_USER := $(shell id -u -n)
POSTGRES_14_TAG ?= 14-alpine

.PHONY: test
test:
	go test ./pkg/...
	