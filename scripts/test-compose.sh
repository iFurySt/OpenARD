#!/usr/bin/env bash
set -euo pipefail

pick_port() {
  python3 - <<'PY'
import socket
with socket.socket() as sock:
    sock.bind(("127.0.0.1", 0))
    print(sock.getsockname()[1])
PY
}

project_name="${ARD_COMPOSE_PROJECT_NAME:-ard-compose-test-$$}"
registry_port="${ARD_COMPOSE_REGISTRY_PORT:-$(pick_port)}"
admin_token="${ARD_COMPOSE_ADMIN_TOKEN:-compose-admin-token}"
registry_url="http://127.0.0.1:${registry_port}"
version="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"
commit="${COMMIT:-$(git rev-parse --short=12 HEAD 2>/dev/null || echo unknown)}"
build_date="${BUILD_DATE:-$(git log -1 --format=%cI 2>/dev/null || date -u +%Y-%m-%dT%H:%M:%SZ)}"

compose() {
  ARD_REGISTRY_PORT="${registry_port}" \
    ARD_ADMIN_TOKEN="${admin_token}" \
    VERSION="${version}" \
    COMMIT="${commit}" \
    BUILD_DATE="${build_date}" \
    docker compose -p "${project_name}" -f infra/compose.yaml "$@"
}

cleanup() {
  compose down -v --remove-orphans >/dev/null 2>&1 || true
}
trap cleanup EXIT

compose down -v --remove-orphans >/dev/null 2>&1 || true
compose up -d --build

health_ready=0
for _ in $(seq 1 60); do
  if curl -fsS "${registry_url}/health" >/dev/null 2>&1; then
    health_ready=1
    break
  fi
  sleep 1
done
if [ "${health_ready}" -ne 1 ]; then
  compose ps || true
  compose logs --no-color registry || true
  exit 1
fi
curl -fsS "${registry_url}/health" >/tmp/ard-compose-health.json
grep -q "\"commit\":\"${commit}\"" /tmp/ard-compose-health.json

bin/ardctl admin add catalog ./internal/catalog/testdata/acme-ai-catalog.json \
  --registry-url "${registry_url}" \
  --admin-token "${admin_token}"
curl -fsS "${registry_url}/.well-known/ai-catalog.json" >/tmp/ard-compose-catalog.json
grep -q "ARD Registry" /tmp/ard-compose-catalog.json
grep -q "Weather Data Node" /tmp/ard-compose-catalog.json
bin/ardctl search weather --registry-url "${registry_url}" --kind mcp --json >/tmp/ard-compose-search.json
grep -q "Weather Data Node" /tmp/ard-compose-search.json

curl -fsS "${registry_url}/metrics" >/tmp/ard-compose-metrics.txt
grep -q "ard_http_requests_total" /tmp/ard-compose-metrics.txt
grep -q "ard_http_request_duration_seconds_bucket" /tmp/ard-compose-metrics.txt
grep -q "ard_runtime_goroutines" /tmp/ard-compose-metrics.txt
