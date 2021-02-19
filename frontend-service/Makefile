windows:
	go build ./cmd/frontend-service/main.go

docker:
	make windows
	docker build -t thijsheijden/microservice-example:latest .
	docker push thijsheijden/microservice-example:latest
	kubectl delete deployment frontend-deployment
	kubectl apply -f ./deployment.yaml
