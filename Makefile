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

.PHONY: skaffold-dev-k3d
skaffold-dev-k3d:
	skaffold dev -p k3d

.PHONY: skaffold-dev-microk8s
skaffold-dev-microk8s:
	skaffold dev -p microk8s

set-local-test:
	$(eval url_value := http://localhost:9080)

set-k3d-test:
	$(eval url_value := http://album-store.local:8070)

.PHONY: local-test
local-test: set-local-test run-tests

.PHONY: k3d-test
k3d-test: set-k3d-test run-tests

run-tests:
	curl --location --request GET '$(url_value)/albums/1' --header 'Accept: application/json';
	curl --location --request GET '$(url_value)/albums/666' --header 'Accept: application/json';
	curl --location --request GET '$(url_value)/albums/X' --header 'Accept: application/json';
	curl --location --request GET '$(url_value)/albums';
	curl --location --request POST '$(url_value)/albums' \
		--header 'Content-Type: application/json' --header 'Accept: application/json' \
		--data-raw '{ "idx": 10, "titlexx": "Blue Train", "artistx": "John Coltrane", "price": 56.99, "X": "asdf" }';
	curl --location --request POST '$(url_value)/albums' \
    		--header 'Content-Type: application/json' --header 'Accept: application/json' \
    		--data-raw '{ "id": -1, "title": "s", "artist": "p", "price": -0.1}';
	curl --location --request POST '$(url_value)/albums' \
        --header 'Content-Type: application/json' --header 'Accept: application/json' \
        --data-raw '{"id": 10,';
	curl --location --request POST '$(url_value)/albums' \
        --header 'Content-Type: application/json' --header 'Accept: application/json' \
        --data-raw '{"id": 10, "title": "The Ozzman Cometh", "artist": "Black Sabbath", "price": 66.60}';
	curl --location --request GET '$(url_value)/v3/api-docs';

.PHONY: docker-build
docker-build:
	docker build -t album-store.local:0.1 -t album-store.local:latest .

.PHONY: k3d-docker-registry
k3d-docker-registry:
	docker tag album-store:latest localhost:54094/album-store:0.1
	docker push localhost:54094/album-store:0.1

.PHONY: k3d-internal-deploy-deployment
k3d-internal-deploy-deployment:
	kubectl apply -f album-store-k3d-deployment.yaml

.PHONY: k3d-internal-undeploy-deployment
k3d-internal-undeploy-deployment:
	kubectl delete -f album-store-k3d-deployment.yaml

.PHONY: k3d-internal-pod
k3d-internal-deploy-pod:
	kubectl apply -f album-store-k3d-pod.yaml

.PHONY: k3d-internal-undeploy-pod
k3d-internal-undeploy-pod:
	kubectl delete -f album-store-k3d-pod.yaml

setup-album-properties:
	$(eval album_setup := GRPC_GO_LOG_SEVERITY_LEVEL=info GRPC_GO_LOG_VERBOSITY_LEVEL=99 NAMESPACE=no-namespace INSTANCE_NAME=album-store-1)

.PHONY: docker-k3d-start
docker-k3d-start: setup-album-properties
	$(album_setup) OTEL_LOCATION=otel-collector.local:8070 docker run -d -p 9080:9080 --name album-store album-store:0.1

.PHONY: docker-local-start
docker-local-start: setup-album-properties
	$(album_setup) OTEL_LOCATION=localhost:4327 docker run -d -p 9080:9080 --name album-store album-store:0.1

.PHONY: local-start-k3d
local-start-k3d: build setup-album-properties
	$(album_setup) OTEL_LOCATION=otel-collector.local:8070 ./album-store

.PHONY: local-start-grpc
local-start-grpc: build setup-album-properties
	$(album_setup) OTEL_LOCATION=localhost:4327 ./album-store

.PHONY: local-start-http
local-start-http: build setup-album-properties
	$(album_setup) OTEL_LOCATION=localhost:4328 ./album-store

.PHONY: docker-stop
docker-stop:
	docker stop album-store;

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
	k3d cluster create k3s-default --config k3d-config.yaml
	kubectl create namespace album-store

.PHONY: k3d-cluster-delete
k3d-cluster-delete:
	k3d cluster delete k3s-default;

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out;
	go tool cover -html=coverage.out -o coverage.html;