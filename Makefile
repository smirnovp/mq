.PHONY: run-publisher
run-publisher:
	go run --race cmd/publisher/publisher.go

.PHONY: run-consumer
run-consumer:
	go run --race cmd/consumer/consumer.go "$(param)"

DEFAULT_GOAL = run-publisher