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

You have to make some modififications.

## Paste the contents 

Add your dashboard json to your `configmap` on the line after the `|-` indented in `4 spaces`.

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
  ...
```

The inputs you cannot use as grafana cannot use a Input dropdown by itself to select a datasource.

We will force the datasouce as described below.

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

```json
"datasource": { 
    "type": "datasource", 
    "uid": "grafana" 
},
```

#### Example annotations datasource block:

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

## Why the replace works

Down the bottom of the dashboard contents in the `templating` section you have where we set the `uid` to be used for Prometheus in this dashboard as a single value.

### Original templating:
```json
"templating": {
    "list": [
      {
        "allValue": null,
        "current": {},
        "datasource": "${DS_PROMETHEUS}",
        ...
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

### Changed templating:
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
        ...
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


## DS_PROMETHEUS change explained 

When:

`"datasource": "${DS_PROMETHEUS}",` 

is replaced with:

```json 
"datasource": {
    "type": "prometheus", 
"uid": "PBFA97CFB590B2093" 
},
``` 

it forces the dashboard to use Prometheus for its datasource.

My prometheus had my datasource as `PBFA97CFB590B2093` so I have kept that value. 

You can use any value as long as it is consistent across the dashboard.

## Common problems.

* I have removed the `__input` fixed the datasource for `grafana` and `prometheus` in the dashboard JSON.
* I have created the configmap and deployed it before the `grafana`
* I have wired in `dashboardproviders.yaml`  
* the name on the left hand side of `dashboardsConfigMaps` matches `dashboardproviders.yaml` name

### Things to try
Reformat your configmap json and add it back in with the right indentation. 

Change the `grafana_dashboard: "1"` to `grafana_dashboard: "2"` if you have 2 dashboards and you only see one at a time. 