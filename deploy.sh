#!/usr/bin/env bash
#
# Deploy Prep Dashboard with PM2.
#   backend  -> Go binary on :9002
#   frontend -> vite preview of built dist on :9001
#
# Usage:  ./deploy.sh
#
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

echo "==> Prep Dashboard deploy"

# --- prerequisites -----------------------------------------------------------
command -v go >/dev/null   || { echo "ERROR: go not found";   exit 1; }
command -v node >/dev/null || { echo "ERROR: node not found"; exit 1; }
command -v npm >/dev/null  || { echo "ERROR: npm not found";  exit 1; }
if ! command -v pm2 >/dev/null; then
  echo "ERROR: pm2 not found. Install it with:  npm install -g pm2"
  exit 1
fi

mkdir -p "$ROOT/logs"

# --- backend: build Go binary ------------------------------------------------
echo "==> Building backend (Go)…"
cd "$ROOT/backend"
go mod download
CGO_ENABLED=0 go build -o prep-backend .
echo "    built backend/prep-backend"

# --- frontend: install deps + production build -------------------------------
echo "==> Building frontend (Vite)…"
cd "$ROOT/frontend"
if [ -f package-lock.json ]; then
  npm ci
else
  npm install
fi
npm run build
echo "    built frontend/dist"

# --- start / reload under PM2 ------------------------------------------------
echo "==> (Re)starting PM2 processes…"
cd "$ROOT"
# startOrReload does a zero-downtime reload if already running, else starts.
pm2 startOrReload ecosystem.config.cjs --update-env
pm2 save

echo ""
pm2 status
echo ""
echo "==> Done."
echo "    Backend : http://localhost:9002/api/health"
echo "    Frontend: http://localhost:9001"
echo ""
echo "    Tip: run 'pm2 startup' once to auto-start on boot."
