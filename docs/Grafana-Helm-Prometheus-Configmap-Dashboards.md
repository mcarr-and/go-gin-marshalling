# Helm Grafana Dashboards setup

## Important note

**Import your `configmap` before you provision your `grafana`**

## Starting point

https://blog.cloudcover.ch/posts/grafana-helm-dashboard-import/

# Grafana Helm Chart Value Changes 

In your Grafana Helm chart you add the following 2 blocks that are commented out.

## dashboardProviders:

In your Helm chart there will be a commented out section 


### Original dashboardProviders: Helm Value

```yaml
dashboardProviders: {}
#  dashboardproviders.yaml:
#    apiVersion: 1
#    providers:
#    - name: 'default'
#      orgId: 1
#      folder: ''
#      type: file
#      disableDeletion: false
#      editable: true
#      options:
#        path: /var/lib/grafana/dashboards/default
```

### Changed dashboardProviders:

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

### What this is doing

Uncommenting the `dashboardproviders.yaml:` and adding a `name` to the `providers` list will allow you to import dashboards.

If you want to add a second dashboard just add it under `providers`

```yaml
    ...
    - name: "golang"
        orgId: 1
        folder: ""
        type: file
        disableDeletion: false
        editable: true
        options:
          path: /var/lib/grafana/dashboards/golang
```

## dashboardsConfigMaps:

### Original dashboardsConfigMaps:

```yaml
dashboardsConfigMaps: {}
#  default: ""
```

### Changed dashboardsConfigMaps:

```yaml
dashboardsConfigMaps:
  jaeger: "dashboard-configmap-jaeger"
```

### What this is doing 

The left hand side of the `:` maps to the name you used in the `providers` block above. E.G.: `jaeger`

The right hand side of the `:` maps to the name of your `configmap` dashboard you want to import from your configmap. E.G.: `dashboard-configmap-jaeger`

# ConfigMap Creation 

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
            ...
```

### Note on Dashboards from Grafana or any plugin you use 

You cannot just download dashboards from https://grafana.com/grafana/dashboards/ and add them fully to your configmap.

You have to make some modifications.

## Import & Export Route for Datasource Changes

You could skip all the step by first import the grafana dashboard into Grafana, select your datasource when you import it.

Then export the dashboard. 

On Export **DO NOT** `Export for sharing externally` as doing this you will have to do all the `datasource` changes that are described below.

You will still have to paste the JSON to the `configmap`

## Grafana Dashboard Modifications Route

Cut and paste your dashboard json to your `configmap` on the line after the `|-` indented in `4 spaces`.

```yaml
...
data:
  jaeger-all-in-one-dashboard.json: |-
    {
      "annotations": {
```

## Remove the inputs 

Remove all the json blocks below. 

```json
  "__inputs": [],
  "__elements": {},
  "__requires": []
  ...
```

The inputs you cannot use as grafana cannot use a Input dropdown by itself to select a datasource.

We will force the `datasource` as described below.

The first json block I have in my dashboards are `"annotations": {`

## Force dashboard datasource to be Grafana

Replace the dropdown value for datasource

`"datasource": "-- Grafana --",` we are going to be removed.

### Original annotations datasource: 

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

### Changed annotations datasource:

The datasource is now: 

`"datasource": "grafana",`

#### Example annotations datasource block:

```json
"annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": "grafana",
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
`"datasource": "prometheus",`

## Why the replacement works

When:

`"datasource": "${DS_PROMETHEUS}",` 

is replaced with:

`"datasource": "prometheus",` 

it forces the dashboard to use prometheus for its datasource.

In my `grafana` Helm values file my `datasources` for `prometheus` has `name: prometheus` so that is why the mapping works.


E.G.: my Grafana Helm Value block for dataources 
```yaml
datasources:
  datasources.yaml:
    apiVersion: 1
    datasources:
      - name: prometheus
        uid: "prometheus"
        type: prometheus
        url: http://prometheus-server.monitoring.svc.cluster.local:80
        access: proxy
        isDefault: true
```

Use the name of your prometheus `datasource` if your Prometheus is using a different `name`.

## Common problems.

### Dashboard not showing up 

* I have removed the `__input` fixed the datasource for `grafana` at the annotations and `prometheus` in the rest of the dashboard JSON.

* I have created the configmap and deployed it before the `grafana`

* I have wired in `dashboardproviders.yaml`  

* the name on the left hand side of `dashboardsConfigMaps` matches `dashboardproviders.yaml` `name`

AND - the dashboard does not show up in Prometheus.

#### Things to try

* Reformat your configmap json and paste the contents back into the `configmap` with the correct 4 space indentation. 

* Change the `grafana_dashboard: "1"` to `grafana_dashboard: "2"` if you have 2 dashboards with the same ID if you only see one or the other. 