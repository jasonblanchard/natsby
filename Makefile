test:
	go test -v -coverprofile=coverage.out

coverage: test
	go tool cover -func=coverage.out

coverage_visualization: test
	go tool cover -html=coverage.out

nats:
	docker run --name nats -d -p 4222:4222 nats

nats_stop:
	docker kill nats
	docker rm nats
