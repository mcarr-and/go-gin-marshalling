.PHONY: build
build:
	go clean;
	go build -ldflags "-X main.version=0.1 -X main.gitHash=`git rev-parse --short HEAD`" -v -o album-store

.PHONY: test
test:
	go test -v

.PHONY: test-benchmark
test-benchmark:
	go test -bench=. -count 2 -run=^# -benchmem

.PHONY: local-start-k3d
local-start-k3d: build
	GRPC_GO_LOG_SEVERITY_LEVEL=info GRPC_GO_LOG_VERBOSITY_LEVEL=99 OTEL_LOCATION=otel-collector.local:8070 INSTANCE_NAME=album-store-1 NAMESPACE=no-namespace ./album-store

.PHONY: local-start-grpc
local-startgrpc:
	GRPC_GO_LOG_SEVERITY_LEVEL=info GRPC_GO_LOG_VERBOSITY_LEVEL=99 OTEL_LOCATION=localhost:4327 INSTANCE_NAME=album-store-1 NAMESPACE=no-namespace ./album-store

.PHONY: local-start-http
local-starthttp: build
	GRPC_GO_LOG_SEVERITY_LEVEL=info GRPC_GO_LOG_VERBOSITY_LEVEL=99 OTEL_LOCATION=localhost:4328 INSTANCE_NAME=album-store-1 NAMESPACE=no-namespace ./album-store

.PHONY: skaffold-dev-k3d
skaffold-dev-k3d:
	skaffold dev -p k3d

.PHONY: skaffold-dev-microk8s
skaffold-dev-microk8s:
	skaffold dev -p microk8s

.PHONY: local-test
local-test:
	curl --location --request GET 'http://localhost:9080/albums/1' --header 'Accept: application/json';
	curl --location --request GET 'http://localhost:9080/albums/666' --header 'Accept: application/json';
	curl --location --request GET 'http://localhost:9080/albums/X' --header 'Accept: application/json';
	curl --location --request GET 'http://localhost:9080/albums';
	curl --location --request POST 'http://localhost:9080/albums' \
		--header 'Content-Type: application/json' --header 'Accept: application/json' \
		--data-raw '{ "idx": 10, "titlexx": "Blue Train", "artistx": "John Coltrane", "price": 56.99, "X": "asdf" }';
	curl --location --request POST 'http://localhost:9080/albums' \
    		--header 'Content-Type: application/json' --header 'Accept: application/json' \
    		--data-raw '{ "id": -1, "title": "s", "artist": "p", "price": -0.1}';
	curl --location --request POST 'http://localhost:9080/albums' \
        --header 'Content-Type: application/json' --header 'Accept: application/json' \
        --data-raw '{"id": 10,';
	curl --location --request POST 'http://localhost:9080/albums' \
        --header 'Content-Type: application/json' --header 'Accept: application/json' \
        --data-raw '{"id": 10, "title": "The Ozzman Cometh", "artist": "Black Sabbath", "price": 66.60}';
	curl --location --request GET 'http://localhost:9080/v3/api-docs';

.PHONY: docker-build
docker-build:
	docker build -t album-store:0.1 -t album-store:latest .

.PHONY: k3d-docker-registry
k3d-docker-registry:
	docker tag album-store:latest localhost:54094/album-store:0.1
	docker push localhost:54094/album-store:0.1

.PHONY: k3d-internal-deploy
k3d-internal-deploy:
	kubectl apply -f album-store-k3d-deployment.yaml

.PHONY: k3d-internal-undeploy
k3d-internal-undeploy:
	kubectl delete -f album-store-k3d-deployment.yaml

.PHONY: docker-k3d-start
docker-k3d-start:
	NAMESPACE=no-namespace INSTANCE_NAME=album-store-1 OTEL_LOCATION=otel-collector.local:8070 docker run -d -p 9080:9080 --name album-store album-store:0.1

.PHONY: docker-local-start
docker-local-start:
	NAMESPACE=no-namespace INSTANCE_NAME=album-store-1 OTEL_LOCATION=localhost:4327 docker run -d -p 9080:9080 --name album-store album-store:0.1

.PHONY: docker-stop
docker-stop:
	docker stop go-gin-example;

.PHONY: docker-compose-full-start
docker-compose-full-start:
	docker-compose up -d;

.PHONY: docker-compose-full-stop
docker-compose-full-stop:
	docker-compose down;

.PHONY: docker-compose-limited-start
docker-compose-limited-start:
	docker-compose -f docker-compose-limited.yaml up -d;

.PHONY: docker-compose-limited-stop
docker-compose-limited-stop:
	docker-compose -f docker-compose-limited.yaml down;

.PHONY: k3d-cluster-create
k3d-cluster-create:
	k3d registry create registry --port 0.0.0.0:54094;
	k3d cluster create k3s-default --config k3d-config.yaml --registry-use k3d-registry:54094;

.PHONY: k3d-cluster-delete
k3d-cluster-delete:
	k3d cluster delete k3s-default;
	k3d registry delete k3d-registry;

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out;
	go tool cover -html=coverage.out -o coverage.html;