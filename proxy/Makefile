.PHONY: build
build:
	go clean;
	go build -ldflags "-X main.version=0.1 -X main.gitHash=`git rev-parse --short HEAD`" -v -o proxy-service

.PHONY: test
test:
	go test -v

#.PHONY: test-benchmark
#test-benchmark:
#	go test -bench=. -count 2 -run=^# -benchmem

set-local-test:
	$(eval url_value := http://localhost:9070)

set-k3d-test:
	$(eval url_value := http://proxy-service.local:8070)

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

.PHONY: docker-build-proxy
docker-build-proxy:
	DOCKER_BUILDKIT=1 docker build -t proxy-service.local:0.1 -t proxy-service.local:latest .

.PHONY: docker-tag-k3d-registry
docker-tag-k3d-registry:
	docker tag proxy-service:latest localhost:54094/proxy-service:0.1
	docker push localhost:54094/proxy-service:0.1

.PHONY: k3d-proxy-deploy-deployment
k3d-proxy-deploy-deployment:
	kubectl apply -f proxy-service-k3d-deployment.yaml

.PHONY: k3d-proxy-undeploy-deployment
k3d-proxy-undeploy-deployment:
	kubectl delete -f proxy-service-k3d-deployment.yaml

.PHONY: k3d-proxy-deploy-pod
k3d-proxy-deploy-pod:
	kubectl apply -f proxy-service-k3d-pod.yaml

.PHONY: k3d-proxy-undeploy-pod
k3d-proxy-undeploy-pod:
	kubectl delete -f proxy-service-k3d-pod.yaml

setup-album-properties:
	$(eval proxy_setup := GRPC_GO_LOG_SEVERITY_LEVEL=info GRPC_GO_LOG_VERBOSITY_LEVEL=99 NAMESPACE=no-namespace INSTANCE_NAME=proxy-service-1)

.PHONY: docker-k3d-start
docker-k3d-start: setup-album-properties
	$(proxy_setup) OTEL_LOCATION=otel-collector.local:8070 ALBUM_STORE_LOCATION=localhost:9080 docker run -d -p 9080:9080 --name proxy-service proxy-service:0.1

.PHONY: docker-local-start
docker-local-start: setup-album-properties
	$(proxy_setup) OTEL_LOCATION=localhost:4327 ALBUM_STORE_LOCATION=localhost:9080 docker run -d -p 9080:9080 --name proxy-service proxy-service:0.1

.PHONY: local-start-k3d
local-start-k3d: build setup-album-properties
	$(proxy_setup) OTEL_LOCATION=otel-collector.local:8070 ALBUM_STORE_LOCATION=album-store.album-store.svc.cluster.local:9080 ./proxy-service

.PHONY: local-start-proxy-grpc
local-start-proxy-grpc: build setup-album-properties
	$(proxy_setup) OTEL_LOCATION=localhost:4327 ALBUM_STORE_LOCATION=localhost:9080 ./proxy-service && cd proxy

.PHONY: local-start-album-grpc
local-start-album-grpc: build setup-album-properties
	cd ../ && $(proxy_setup) OTEL_LOCATION=localhost:4327 ./album-store && cd proxy

.PHONY: local-start-proxy-http
local-start-proxy-http: build setup-album-properties
	$(proxy_setup) OTEL_LOCATION=localhost:4328 ALBUM_STORE_LOCATION=localhost:9080 ./proxy-service

.PHONY: docker-stop
docker-stop:
	docker stop proxy-service;

.PHONY: docker-compose-limited-start
docker-compose-limited-start:
	cd .. && $(MAKE) docker-compose-limited-start && cd proxy

.PHONY: skaffold-dev-k3d
skaffold-dev-k3d:
	cd .. && $(MAKE) skaffold-dev-k3d && cd proxy

.PHONY: k3d-cluster-create
k3d-cluster-create:
	cd .. && $(MAKE) k3d-cluster-create && cd proxy

.PHONY: k3d-album-deploy-deployment
k3d-album-deploy-deployment:
	cd .. && $(MAKE) k3d-album-deploy-deployment && cd proxy

.PHONY: k3d-album-undeploy-deployment
k3d-album-undeploy-deployment:
	cd .. && $(MAKE) k3d-album-undeploy-deployment && cd proxy

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out;
	go tool cover -html=coverage.out -o coverage.html;