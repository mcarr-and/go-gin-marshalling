
https://kubernetes.github.io/ingress-nginx/troubleshooting/

## Get Nginx configuration
```bash
kubectl exec -it -n ingress-nginx $(kubectl -n ingress-nginx get pods -l app.kubernetes.io/name=ingress-nginx -o jsonpath="{.items[0].metadata.name}") -- cat /etc/nginx/nginx.conf > nginx.conf
```

## Get Nginx Logs
```bash
kubectl logs -f -n ingress-nginx $(kubectl -n ingress-nginx get pods -l app.kubernetes.io/name=ingress-nginx -o jsonpath="{.items[0].metadata.name}")
```


http://album-store.local:8070/albums/

Gives a 503

------------ [error] 241#241: *720 connect() failed (111: Connection refused) while connecting to upstream, client: 10.42.0.0, server: album-store.local, request: "GET /albums/ HTTP/1.1", upstream: "http://10.42.0.104:9080/albums/", host: "album-store.local:8070"
10.42.0.0 - - [---------- +0000] "GET /albums/ HTTP/1.1" 502 552 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36" 476 0.003 [opentelemetry-album-store-9080] [] 10.42.0.104:9080, 10.42.0.104:9080, 10.42.0.104:9080 0, 0, 0 0.001, 0.001, 0.001 502, 502, 502 a83c834989b5144500d276644fea9cc3

```bash
curl -v http://album-store.local:8070/albums/ GET -H "Content-Type: application/json" -H "Host: http://album-store.local:8070
```

Gives a 400 bad request

10.42.0.0 - - [-----] "GET /albums/ HTTP/1.1" 502 150 "-" "curl/7.84.0" 93 0.000 [opentelemetry-album-store-9080] [] 10.42.1.12:9080, 10.42.1.12:9080, 10.42.1.12:9080 0, 0, 0 0.000, 0.001, 0.000 502, 502, 502 57f76a972cc8be95a9d796847a264eea
----- [error] 242#242: *1930 connect() failed (111: Connection refused) while connecting to upstream, client: 10.42.0.0, server: album-store.local, request: "GET /albums/ HTTP/1.1", upstream: "http://10.42.1.12:9080/albums/", host: "album-store.local:8070"
----- [error] 242#242: *1930 connect() failed (111: Connection refused) while connecting to upstream, client: 10.42.0.0, server: album-store.local, request: "GET /albums/ HTTP/1.1", upstream: "http://10.42.1.12:9080/albums/", host: "album-store.local:8070"
---- [error] 242#242: *1930 connect() failed (111: Connection refused) while connecting to upstream, client: 10.42.0.0, server: album-store.local, request: "GET /albums/ HTTP/1.1", upstream: "http://10.42.1.12:9080/albums/", host: "album-store.local:8070"

Internal issue with K3d and all the minikube type provblems.

https://stackoverflow.com/questions/62162209/ingress-nginx-errors-connection-refused

https://stackoverflow.com/questions/72097715/kubernetes-k3d-ingress-always-error-404-not-found

## Get Album Store Logs 

```bash
kubectl logs -f -n opentelemetry $(kubectl -n opentelemetry get pods -l app.kubernetes.io/name=album-store -o jsonpath="{.items[0].metadata.name}")
```

## Get opentelemetry Logs

```bash
kubectl logs -f -n opentelemetry $(kubectl -n opentelemetry get pods  -l app.kubernetes.io/name=opentelemetry-collector -o jsonpath="{.items[0].metadata.name}")
```
