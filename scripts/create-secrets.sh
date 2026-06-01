#!/usr/bin/env bash
#
# create-secrets.sh — interactively create the Kubernetes secrets that
# api/deployment.yaml references (project.md §9.6 — secrets via secretKeyRef only).
#
# Secret created:  naharda-api  (namespace: naharda)
#   DATABASE_URL        (required)  — Postgres DSN, e.g.
#                                     postgres://user:pass@host:5432/naharda?sslmode=require
#   ALERT_WEBHOOK_URL   (optional)  — data-quality alert webhook (ntfy/Slack)
#
# web/deployment.yaml needs NO secret (its env is inline).
#
# Requires: kubectl with a context pointing at the target cluster.
set -euo pipefail

NS="${1:-naharda}"
SECRET="naharda-api"

command -v kubectl >/dev/null 2>&1 || { echo "kubectl not found in PATH"; exit 1; }

echo "Cluster context: $(kubectl config current-context)"
echo "Namespace:       $NS"
echo "Secret:          $SECRET"
echo
read -r -p "Proceed against this cluster/namespace? [y/N] " ok
[[ "$ok" =~ ^[Yy]$ ]] || { echo "aborted"; exit 1; }

# Ensure namespace exists.
kubectl get namespace "$NS" >/dev/null 2>&1 || {
  echo "Creating namespace $NS"
  kubectl create namespace "$NS"
}

# DATABASE_URL (required, hidden input).
read -r -s -p "DATABASE_URL (Postgres DSN): " DATABASE_URL; echo
[[ -n "$DATABASE_URL" ]] || { echo "DATABASE_URL is required"; exit 1; }

# ALERT_WEBHOOK_URL (optional).
read -r -p "ALERT_WEBHOOK_URL (optional, blank to skip): " ALERT_WEBHOOK_URL

args=(--from-literal=DATABASE_URL="$DATABASE_URL")
if [[ -n "${ALERT_WEBHOOK_URL:-}" ]]; then
  args+=(--from-literal=ALERT_WEBHOOK_URL="$ALERT_WEBHOOK_URL")
fi

# Idempotent upsert (create-or-replace) without echoing secret values.
kubectl create secret generic "$SECRET" -n "$NS" "${args[@]}" \
  --dry-run=client -o yaml | kubectl apply -f -

echo
echo "✓ Secret '$SECRET' applied in namespace '$NS'."
echo "  Keys: DATABASE_URL${ALERT_WEBHOOK_URL:+, ALERT_WEBHOOK_URL}"
