#!/usr/bin/env bash
set -euo pipefail

METADATA_URL="http://metadata.google.internal/computeMetadata/v1/instance/attributes"

metadata_get() {
  local key="$1"
  curl -fsS -H "Metadata-Flavor: Google" "${METADATA_URL}/${key}" || true
}

APP_NAME="$(metadata_get APP_NAME)"
IMAGE="$(metadata_get IMAGE)"
ENV_FILE="/etc/wpokt-validator.env"

if [[ -z "${APP_NAME}" ]]; then
  APP_NAME="wpokt-validator"
fi

if [[ -z "${IMAGE}" ]]; then
  IMAGE="docker.io/dan13ram/wpokt-validator:v0.2.7"
fi

if ! command -v docker >/dev/null 2>&1; then
  apt-get update
  apt-get install -y docker.io
  systemctl enable --now docker
fi

mkdir -p /etc
umask 077

MONGODB_DATABASE="$(metadata_get MONGODB_DATABASE)"
ETH_RPC_URL="$(metadata_get ETH_RPC_URL)"
POKT_RPC_URL="$(metadata_get POKT_RPC_URL)"
GOOGLE_SECRET_MANAGER_ENABLED="$(metadata_get GOOGLE_SECRET_MANAGER_ENABLED)"
GOOGLE_MONGO_SECRET_NAME="$(metadata_get GOOGLE_MONGO_SECRET_NAME)"
GOOGLE_POKT_SECRET_NAME="$(metadata_get GOOGLE_POKT_SECRET_NAME)"
GOOGLE_ETH_SECRET_NAME="$(metadata_get GOOGLE_ETH_SECRET_NAME)"
POKT_GCP_KMS_KEY_NAME="$(metadata_get POKT_GCP_KMS_KEY_NAME)"
POKT_START_HEIGHT="$(metadata_get POKT_START_HEIGHT)"
ETH_START_BLOCK_NUMBER="$(metadata_get ETH_START_BLOCK_NUMBER)"
LOG_LEVEL="$(metadata_get LOG_LEVEL)"

cat > "${ENV_FILE}" <<EOF
MONGODB_DATABASE=${MONGODB_DATABASE}
ETH_RPC_URL=${ETH_RPC_URL}
POKT_RPC_URL=${POKT_RPC_URL}
GOOGLE_SECRET_MANAGER_ENABLED=${GOOGLE_SECRET_MANAGER_ENABLED}
GOOGLE_MONGO_SECRET_NAME=${GOOGLE_MONGO_SECRET_NAME}
GOOGLE_POKT_SECRET_NAME=${GOOGLE_POKT_SECRET_NAME}
GOOGLE_ETH_SECRET_NAME=${GOOGLE_ETH_SECRET_NAME}
POKT_GCP_KMS_KEY_NAME=${POKT_GCP_KMS_KEY_NAME}
POKT_START_HEIGHT=${POKT_START_HEIGHT}
ETH_START_BLOCK_NUMBER=${ETH_START_BLOCK_NUMBER}
LOG_LEVEL=${LOG_LEVEL}
EOF

docker pull "${IMAGE}"
docker rm -f "${APP_NAME}" || true
docker run -d \
  --name "${APP_NAME}" \
  --restart=always \
  --env-file "${ENV_FILE}" \
  "${IMAGE}"
