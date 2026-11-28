#!/usr/bin/env sh
# Dijalankan DI SERVER aaPanel (bukan di laptop) setelah repo di-clone.
# Menarik kode terbaru dari GitHub, build binary Linux di server, lalu restart daemon.
# File .env tidak diatur di sini: ia rahasia dan dibuat manual sekali di folder backend/.
set -e

APP_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(cd "$APP_DIR/.." && pwd)"
DAEMON_NAME="${DAEMON_NAME:-siminventaris-api}"

echo "==> Menarik kode terbaru di $REPO_DIR"
cd "$REPO_DIR"
git pull --ff-only

echo "==> Build binary di $APP_DIR"
cd "$APP_DIR"
CGO_ENABLED=0 go build -o siminventaris-api ./cmd/api

echo "==> Restart daemon $DAEMON_NAME"
if command -v supervisorctl >/dev/null 2>&1; then
  supervisorctl restart "$DAEMON_NAME"
else
  echo "supervisorctl tidak ditemukan: restart daemon $DAEMON_NAME lewat aaPanel (Go Project)."
fi

echo "==> Selesai. Cek: curl http://127.0.0.1:3000/api/status"
