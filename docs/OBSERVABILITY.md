# Observability

This document gives operators and agents a local workflow for checking registry health,
metrics, logs, and traces.

## Health

```sh
ardctl health --registry-url http://127.0.0.1:8080 --json
```

`GET /health` returns `status`, active entry count, and build metadata.

## Metrics

```sh
ardctl metrics --registry-url http://127.0.0.1:8080
```

`GET /metrics` returns Prometheus text with registry uptime, in-flight requests, request
totals, HTTP duration histograms, goroutines, heap, and GC state.

## Logs

Every HTTP request emits one JSON access log line with `requestId`, `traceId`, `spanId`,
method, path, status, latency, and client IP. Logs do not include bearer tokens or
request bodies.

## Traces

Trace export is disabled by default. Enable OTLP/HTTP trace export by pointing the server
at an OpenTelemetry collector traces endpoint:

```sh
ARD_OTLP_TRACES_ENDPOINT=http://127.0.0.1:4318/v1/traces \
ard-server --database-url "$DATABASE_URL" --addr :8080
```

The shorthand base URL also works:

```sh
ard-server --otlp-traces-endpoint http://127.0.0.1:4318
```

When enabled, the registry exports one server span per HTTP request with the W3C trace
ID/span ID, parent span ID when an inbound `traceparent` is present, request ID, method,
path, route template, status, latency timestamps, and client address.

The local real-artifact E2E gate starts a temporary OTLP/HTTP capture endpoint and
verifies that a live registry request exports a trace span:

```sh
make test-e2e
```

## Current Gaps

- Sampling is all-or-nothing through the exporter endpoint toggle.
- Dashboards and alert rules are deployment-owned and not included yet.
