.PHONY: build
build:
	go mod tidy;
	go get;
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

.PHONY: eval-git-hash
eval-git-hash:
	$(eval GIT_HASH:= $(shell git rev-parse --short HEAD))
	echo $(GIT_HASH)


.PHONY: docker-build-album
docker-build-album: eval-git-hash
	DOCKER_BUILDKIT=1 docker build --build-arg GIT_HASH=$(GIT_HASH) -t album-store:0.1 -t album-store:latest .

.PHONY: docker-tag-k3d-registry-album
docker-tag-k3d-registry-album: docker-build-album
	docker tag album-store:latest localhost:54094/album-store:latest
	docker tag album-store:0.1 localhost:54094/album-store:0.1
	docker push localhost:54094/album-store:latest
	docker push localhost:54094/album-store:0.1

.PHONY: docker-build-proxy
docker-build-proxy:
	cd proxy && $(MAKE) docker-build-proxy && cd ..

.PHONY: docker-tag-k3d-registry-proxy
docker-tag-k3d-registry-proxy: docker-build-proxy
	cd proxy && $(MAKE) docker-tag-k3d-registry-proxy && cd ..

.PHONY: k3d-album-deploy-deployment
k3d-album-deploy-deployment: docker-tag-k3d-registry-album
	kubectl apply -f ./install/album-store-k3d-deployment.yaml

.PHONY: k3d-album-undeploy-deployment
k3d-album-undeploy-deployment:
	kubectl delete -f ./install/album-store-k3d-deployment.yaml

.PHONY: k3d-album-deploy-pod
k3d-album-deploy-pod: docker-tag-k3d-registry-album
	kubectl apply -f ./install/album-store-k3d-pod.yaml

.PHONY: k3d-album-undeploy-pod
k3d-album-undeploy-pod:
	kubectl delete -f ./install/album-store-k3d-pod.yaml

.PHONY: k3d-proxy-deploy-deployment
k3d-proxy-deploy-deployment: docker-tag-k3d-registry-proxy
	kubectl apply -f ./install/proxy-service-k3d-deployment.yaml

.PHONY: k3d-proxy-undeploy-deployment
k3d-proxy-undeploy-deployment:
	kubectl delete -f ./install/proxy-service-k3d-deployment.yaml

.PHONY: k3d-proxy-deploy-pod
k3d-proxy-deploy-pod: docker-tag-k3d-registry-proxy
	kubectl apply -f ./install/proxy-service-k3d-pod.yaml

.PHONY: k3d-proxy-undeploy-pod
k3d-proxy-undeploy-pod:
	kubectl delete -f proxy/proxy-service-k3d-pod.yaml

setup-album-properties:
	$(eval album_setup := GRPC_GO_LOG_SEVERITY_LEVEL=info GRPC_GO_LOG_VERBOSITY_LEVEL=99 NAMESPACE=no-namespace INSTANCE_NAME=album-store-1)

setup-album-docker-properties:
	$(eval album_setup := -e GRPC_GO_LOG_SEVERITY_LEVEL=info -e GRPC_GO_LOG_VERBOSITY_LEVEL=99 -e NAMESPACE=no-namespace -e INSTANCE_NAME=album-store-1)

.PHONY: docker-k3d-start
docker-k3d-start: docker-build-album setup-album-docker-properties
	docker run -d -p 9080:9080 $(album_setup) -e OTEL_LOCATION=otel-collector.local:8070 --name album-store album-store:0.1

.PHONY: docker-local-start
docker-local-start: docker-build-album setup-album-docker-properties
	docker run -d -p 9080:9080  $(album_setup) -e OTEL_LOCATION=localhost:4327 --name album-store album-store:0.1

.PHONY: local-start-k3d
local-start-k3d: build setup-album-properties
	$(album_setup) OTEL_LOCATION=otel-collector.local:8070 ./album-store

.PHONY: local-start-grpc
local-start-grpc: build setup-album-properties
	$(album_setup) OTEL_LOCATION=localhost:4327 ./album-store

.PHONY: docker-stop
docker-stop:
	docker stop album-store;

.PHONY: docker-compose-full-start
docker-compose-full-start: docker-build-album
	docker-compose -f ./install/docker-compose.yaml -d up --remove-orphans;

.PHONY: docker-compose-full-stop
docker-compose-full-stop:
	docker-compose -f ./install/docker-compose.yaml down;

.PHONY: docker-compose-limited-start
docker-compose-limited-start:
	docker-compose -f ./install/docker-compose-limited.yaml up -d --remove-orphans;

.PHONY: docker-compose-limited-stop
docker-compose-limited-stop:
	docker-compose -f ./install/docker-compose-limited.yaml down;

.PHONY: k3d-cluster-create
k3d-cluster-create:
	k3d cluster create k3s-default --config ./instal/k3d-config.yaml
	kubectl create namespace album-store
	kubectl create namespace proxy-service

.PHONY: k3d-cluster-delete
k3d-cluster-delete:
	k3d cluster delete k3s-default;

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out;
	go tool cover -html=coverage.out -o coverage.html;