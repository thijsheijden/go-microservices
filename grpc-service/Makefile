greeting-service:
	go build ./cmd/grpc-service/main.go
	docker build -t thijsheijden/greeter-service:latest .
	docker push thijsheijden/greeter-service:latest