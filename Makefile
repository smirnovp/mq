.PHONY: run-publisher
run-publisher:
	go run cmd/publisher/publisher.go

.PHONY: run-consumer
run-consumer:
	go run cmd/consumer/consumer.go

DEFAULT_GOAL := run-publisher
