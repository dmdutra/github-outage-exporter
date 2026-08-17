# GitHub Outage Exporter

Prometheus exporter that collects outage and degradation metrics from [GitHub Status](https://www.githubstatus.com/) via the [official API](https://www.githubstatus.com/api#status).

## Consumed endpoints

| Endpoint | Description |
|----------|-------------|
| `/api/v2/status.json` | Overall platform status |
| `/api/v2/components.json` | Per-component status |
| `/api/v2/incidents/unresolved.json` | Open incidents |
| `/api/v2/scheduled-maintenances/active.json` | Active scheduled maintenances |

## Exposed metrics

| Metric | Type | Description |
|--------|------|-------------|
| `github_status_indicator_severity` | gauge | Overall severity (`0=none`, `1=minor`, `2=major`, `3=critical`) |
| `github_component_status_severity` | gauge | Per-component severity (`0=operational`, `1=degraded`, `2=partial_outage`, `3=major_outage`) |
| `github_component_operational` | gauge | `1` if the component is operational, `0` otherwise |
| `github_incident_active` | gauge | `1` for each unresolved incident |
| `github_maintenance_active` | gauge | `1` for each active scheduled maintenance |
| `github_scrape_success` | gauge | `1` if the last scrape succeeded |
| `github_scrape_duration_seconds` | gauge | Duration of the last scrape |

## Local execution

```bash
go run ./cmd/github-outage-exporter
```

Visit `http://localhost:8080/metrics`.

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LISTEN_ADDR` | `:8080` | HTTP listen address |
| `GITHUB_STATUS_URL` | `https://www.githubstatus.com` | Status page base URL |
| `SCRAPE_TIMEOUT` | `10s` | Timeout per API request |
| `METRICS_PATH` | `/metrics` | Prometheus metrics path |

## Docker

```bash
docker build -t github-outage-exporter .
docker run --rm -p 8080:8080 github-outage-exporter
```

## Prometheus scrape example

```yaml
scrape_configs:
  - job_name: github-outage-exporter
    static_configs:
      - targets: ["localhost:8080"]
```

## Alert example

```yaml
groups:
  - name: github-status
    rules:
      - alert: GitHubComponentDegraded
        expr: github_component_operational == 0
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "GitHub component degraded: {{ $labels.component }}"
          description: "Component {{ $labels.component }} has status {{ $labels.status }}."

      - alert: GitHubIncidentActive
        expr: github_incident_active == 1
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Active GitHub incident: {{ $labels.name }}"
          description: "Impact {{ $labels.impact }}, status {{ $labels.status }}."
```
