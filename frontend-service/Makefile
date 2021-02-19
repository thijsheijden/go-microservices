windows:
	go build ./cmd/frontend-service/main.go

docker:
	make windows
	docker build -t thijsheijden/microservice-example:latest .
	docker push thijsheijden/microservice-example:latest