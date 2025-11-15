#!/usr/bin/env bash
set -euo pipefail

# generate_certs.sh
# Generates a self-signed certificate and private key for local HTTPS testing.
# Output files: certs/server.crt and certs/server.key (overwrites existing files).

OUT_DIR="$(cd "$(dirname "$0")" && pwd)"
CRT="$OUT_DIR/server.crt"
KEY="$OUT_DIR/server.key"

COMMON_NAME="localhost"
DAYS=365

echo "Generating self-signed certificate for ${COMMON_NAME} into ${OUT_DIR}"

openssl req -x509 -nodes -days ${DAYS} -newkey rsa:2048 \
  -keyout "$KEY" -out "$CRT" \
  -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=${COMMON_NAME}" \
  -addext "subjectAltName = DNS:localhost, IP:127.0.0.1"

chmod 644 "$CRT"
chmod 600 "$KEY"

echo "Wrote: $CRT"
echo "Wrote: $KEY"
echo "If you plan to commit certs for demo only, change .gitignore in certs/ accordingly."
