#!/bin/bash
# Build script for PostgreSQL 16 + Apache AGE + pgvector Docker image
# Use this when `docker compose build postgres` fails due to DNS issues.
#
# Usage: bash scripts/build-postgres-image.sh

set -euo pipefail

IMAGE_NAME="jolly-eagle-postgres"
CONTAINER_NAME="pg-build-$$"
AGE_BRANCH="PG16/v1.5.0-rc0"

echo "==> Starting build container with DNS override..."
docker run --dns 1.1.1.1 --name "$CONTAINER_NAME" -d postgres:16-bookworm tail -f /dev/null

cleanup() {
    echo "==> Cleaning up..."
    docker stop "$CONTAINER_NAME" 2>/dev/null || true
    docker rm -f "$CONTAINER_NAME" 2>/dev/null || true
}
trap cleanup EXIT

echo "==> Installing build dependencies and pgvector..."
docker exec "$CONTAINER_NAME" bash -c "
    export DEBIAN_FRONTEND=noninteractive
    apt-get update && apt-get install -y --no-install-recommends \
        build-essential git postgresql-server-dev-16 \
        libreadline-dev zlib1g-dev flex bison \
        postgresql-16-pgvector \
    && rm -rf /var/lib/apt/lists/*
"

echo "==> Building Apache AGE ($AGE_BRANCH)..."
docker exec "$CONTAINER_NAME" bash -c "
    cd /tmp && \
    git -c http.sslVerify=false clone --branch $AGE_BRANCH --depth 1 https://github.com/apache/age.git && \
    cd age && \
    make PG_CONFIG=/usr/lib/postgresql/16/bin/pg_config && \
    make install PG_CONFIG=/usr/lib/postgresql/16/bin/pg_config
"

echo "==> Configuring PostgreSQL..."
docker exec "$CONTAINER_NAME" bash -c "
    echo \"shared_preload_libraries = 'age'\" >> /usr/share/postgresql/postgresql.conf.sample && \
    mkdir -p /docker-entrypoint-initdb.d && \
    echo 'CREATE EXTENSION IF NOT EXISTS vector;' > /docker-entrypoint-initdb.d/00-init-extensions.sql && \
    echo 'CREATE EXTENSION IF NOT EXISTS age;' >> /docker-entrypoint-initdb.d/00-init-extensions.sql && \
    echo \"LOAD 'age';\" >> /docker-entrypoint-initdb.d/00-init-extensions.sql && \
    echo \"SET search_path = ag_catalog, \\\"\\\$user\\\", public;\" >> /docker-entrypoint-initdb.d/00-init-extensions.sql && \
    rm -rf /tmp/age
"

echo "==> Stopping container..."
docker stop "$CONTAINER_NAME"

echo "==> Exporting clean image (without data directory)..."
# Export the filesystem, then import it — this strips the data volume state
TEMP_TAR=$(mktemp /tmp/pg-export.XXXXXX.tar)
docker export "$CONTAINER_NAME" > "$TEMP_TAR"
docker import "$TEMP_TAR" "$IMAGE_NAME"
rm -f "$TEMP_TAR"

echo "==> Tagging image..."
docker tag "$IMAGE_NAME" "$IMAGE_NAME:latest"

echo ""
echo "==> Build complete! Image: $IMAGE_NAME"
echo "==> Verify extensions with:"
echo "    docker compose up -d postgres"
echo "    docker exec graph-rag-postgres psql -U graphrag -d graphrag -c 'SELECT extname, extversion FROM pg_extension;'"
