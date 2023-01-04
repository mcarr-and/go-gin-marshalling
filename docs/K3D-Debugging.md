
https://kubernetes.github.io/ingress-nginx/troubleshooting/

## Get Nginx configuration
```bash
kubectl exec -it -n ingress-nginx $(kubectl -n ingress-nginx get pods -l app.kubernetes.io/name=ingress-nginx -o jsonpath="{.items[0].metadata.name}") -- cat /etc/nginx/nginx.conf > nginx.conf
```

## Get Nginx Logs
```bash
kubectl logs -f -n ingress-nginx $(kubectl -n ingress-nginx get pods -l app.kubernetes.io/name=ingress-nginx -o jsonpath="{.items[0].metadata.name}")
```

```bash
curl --insecure --location 'http://album-store.local:8070/albums/'; 
```

## Get Album Store Logs 

```bash
kubectl logs -f -n opentelemetry $(kubectl -n opentelemetry get pods -l app.kubernetes.io/name=album-store -o jsonpath="{.items[0].metadata.name}")
```

## Get opentelemetry Logs

```bash
kubectl logs -f -n opentelemetry $(kubectl -n opentelemetry get pods  -l app.kubernetes.io/name=opentelemetry-collector -o jsonpath="{.items[0].metadata.name}")
```
