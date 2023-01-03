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
	export GRPC_GO_LOG_VERBOSITY_LEVEL=99;
	export GRPC_GO_LOG_SEVERITY_LEVEL=info;
	export NAMESPACE=no-namespace;
	export INSTANCE_NAME=album-store-1;
	export OTEL_LOCATION=otel-collector.local:8070;
	./album-store

.PHONY: local-start-docker-compose-grpc
local-start-docker-compose-grpc: build
	export GRPC_GO_LOG_VERBOSITY_LEVEL=99;
	export GRPC_GO_LOG_SEVERITY_LEVEL=info;
	export NAMESPACE=no-namespace;
	export INSTANCE_NAME=album-store-1;
	export OTEL_LOCATION=otel-collector.local:8070;
	./album-store -otel-location=localhost:4327 -namespace=no-namespace -instance-name=album-store-1

.PHONY: local-start-docker-compose-http
local-start-docker-compose-http: build
	export NAMESPACE=no-namespace;
	export INSTANCE_NAME=album-store-1;
	export OTEL_LOCATION=otel-collector.local:8070;
	./album-store -otel-location=localhost:4328 -namespace=no-namespace -instance-name=album-store-1

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
	docker build -t album-store:0.1 -t album-store .

.PHONY: docker-build-k3d
docker-build-k3d:
	export DOCKER_HOST="tcp://localhost:54094";
	docker build -t album-store:0.1 -t album-store .
	docker tag album-store:latest localhost:54094/album-store:v0.1
	docker push localhost:54094/album-store:v0.1

.PHONY: docker-start
docker-start:
	export NAMESPACE=no-namespace;
	export INSTANCE_NAME=album-store-1;
	export OTEL_LOCATION=otel-collector.local:8070;
	docker run -d -p 9080:9080 --name go-gin-example go-gin-example:0.1

.PHONY: docker-stop
docker-stop:
	docker stop go-gin-example;

.PHONY: docker-compose-start
docker-compose-start:
	docker-compose up -d;

.PHONY: docker-compose-stop
docker-compose-stop:
	docker-compose down;

.PHONY: k3d-cluster-create
k3d-cluster-create:
	k3d cluster create k3s-default --config k3d-config.yaml;
	k3d registry create --port 0.0.0.0:54094;

.PHONY: k3d-cluster-delete
k3d-cluster-delete:
	k3d cluster delete k3s-default;

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out;
	go tool cover -html=coverage.out -o coverage.html;