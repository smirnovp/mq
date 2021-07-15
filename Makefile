.PHONY: run-publisher
run-publisher:
	go run cmd/publisher/publisher.go 

.PHONY: run-consumer
run-consumer:
	go run cmd/consumer/consumer.go "$(param)"

DEFAULT_GOAL := run-publisher
