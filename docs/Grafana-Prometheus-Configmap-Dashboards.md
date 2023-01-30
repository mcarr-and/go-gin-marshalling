Starting point
https://blog.cloudcover.ch/posts/grafana-helm-dashboard-import/



In your Grafana Helm chart you add the following 2 blocks that are commented out

```yaml
dashboardProviders:
  dashboardproviders.yaml:
    apiVersion: 1
    providers:
      - name: "jaeger"
        orgId: 1
        folder: ""
        type: file
        disableDeletion: false
        editable: true
        options:
          path: /var/lib/grafana/dashboards/jaeger
```

AND something like the following under `dashboardsConfigMaps` 



```yaml
dashboardsConfigMaps:
  jaeger: "dashboard-configmap-jaeger"
```


The left hand side of the `:` maps to the name you used in the `providers` block above. E.G.: `jaeger`

The right hand side of the `:` maps to the dashboard you want to import from your configmap. E.G.: `dashboard-configmap-jaeger`


## ConfigMap 

Example of a `dashboard-configmap-jaeger`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    grafana_dashboard: "1"
  name: dashboard-configmap-jaeger
  namespace: monitoring
data:
  jaeger-all-in-one-dashboard.json: |-
    {
      "annotations": {
        "list": [
          {
            "builtIn": 1,
            "datasource": {
              "type": "datasource",
              "uid": "grafana"
            },
            ....
```

### Note on Dashboards from Grafana or any Plugin you use 

You cannot just download dashboards from https://grafana.com/grafana/dashboards/ and add them fully to your configmap.

You have to make some modififications.

## Paste the contents 

Add your dashboard json on the line after the `|-` indented in `4 spaces`.

```yaml
...

data:
  jaeger-all-in-one-dashboard.json: |-
    {
      "annotations": {
```


## Remove the inputs 

Remove all the blocks below. 



```json
  "__inputs": [],
  "__elements": {},
  "__requires": []
```

The inputs you cannot use as you want to force the datasource as grafana cannot use a Input dropdown by itself.

We will force the datasouce as described below.


## Force datasource to be Grafana

### Replace the dropdown value 

`"datasource": "-- Grafana --",` we are going to remove 

#### Original value 
```json
"annotations": {
    "list": [
        {
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
        }
    ]
},
```

#### Changed datasource

`"datasource": { "type": "datasource", "uid": "grafana" },` is now set.

```json
"annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": {
          "type": "datasource",
          "uid": "grafana"
        },
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
    ]
},
```

## Replace all instances of ${DS_PROMETHEUS} 

### Find all instances of ${DS_PROMETHEUS} 

```json
"datasource": {
  "type": "prometheus",
  "uid": "${DS_PROMETHEUS}"
},
```

### Replace with fixed datasource
```json
"datasource": {
  "type": "prometheus",
  "uid": "PBFA97CFB590B2093"
},
```

### Why the replace works

Down the bottom of the dashboard contents in the `templating` section you have where we set the `uid` to be used for Prometheus in this dashboard as a single value.

#### Original templating
```json
"templating": {
    "list": [
      {
        "allValue": null,
        "current": {},
        "datasource": "${DS_PROMETHEUS}",
        .....
      },
      {
        "allValue": null,
        "current": {},
        "datasource": "${DS_PROMETHEUS}",
        ...
      }
    ]
  },
```

#### Forced templating
```json
"templating": {
    "list": [
      {
        "allValue": null,
        "current": {},
        "datasource": {
          "type": "prometheus",
          "uid": "PBFA97CFB590B2093"
        },
        .....
      },
      {
        "allValue": null,
        "current": {},
        "datasource": {
          "type": "prometheus",
          "uid": "PBFA97CFB590B2093"
        },
        ...
      }
    ]
  },
```

When 

`"datasource": "${DS_PROMETHEUS}",` 

is replaced with the 

`"datasource": {"type": "prometheus", "uid": "PBFA97CFB590B2093" },` 

it forces this dashboard to use Prometheus for its datasource.

My prometheus had my datasource as `PBFA97CFB590B2093` so I have kept that value. 

You can use any value as long as it is consistent across the dashboard.
