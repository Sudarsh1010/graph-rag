# Dokploy Deployment Guide

This guide covers deploying the Graph RAG application to a VPS using Dokploy with Nixpacks.

## Prerequisites

- Dokploy installed on your VPS ([Installation Guide](https://docs.dokploy.com))
- A Git repository with this code
- Domain configured (optional, but recommended for SSL)

## Quick Start

### 1. Create PostgreSQL Database

1. In Dokploy, navigate to your **Project** → **Create Service** → **Database**
2. Select **PostgreSQL**
3. Configure:
   - **Name**: `graph-rag-db` (this becomes the hostname)
   - **Database Name**: `graphrag`
   - **Database User**: `graphrag`
   - **Database Password**: `<your-secure-password>`
   - **Docker Image**: `postgres:16` (or leave default)
4. Leave **External Port** empty (internal-only access for security)
5. Click **Deploy**

### 2. Deploy the Application

1. In the same project, **Create Service** → **Application**
2. Connect your Git repository
3. Configure:
   - **Build Type**: Nixpacks (default)
   - **Branch**: `main` (or your preferred branch)
4. Add **Environment Variables**:

```bash
POSTGRES_HOST=graph-rag-db
POSTGRES_PORT=5432
POSTGRES_USERNAME=graphrag
POSTGRES_PASSWORD=<your-secure-password>
POSTGRES_DB=graphrag
POSTGRES_SSLMODE=disable
LOG_LEVEL=info
LOG_FORMAT=json
ENV=production
PORT=8080
```

5. Click **Deploy**

### 3. Configure Domain (Optional)

1. Go to **Application** → **Domains** → **Add Domain**
2. Enter your domain (e.g., `api.yourdomain.com`)
3. Dokploy provisions SSL automatically via Traefik

## Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `POSTGRES_HOST` | PostgreSQL hostname | - | Yes |
| `POSTGRES_PORT` | PostgreSQL port | `5432` | Yes |
| `POSTGRES_USERNAME` | Database username | - | Yes |
| `POSTGRES_PASSWORD` | Database password | - | Yes |
| `POSTGRES_DB` | Database name | - | Yes |
| `POSTGRES_SSLMODE` | SSL mode for connection | `disable` | No |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` | No |
| `LOG_FORMAT` | Log format (console, json) | `console` | No |
| `ENV` | Environment (development, production) | `development` | No |
| `PORT` | Server port | `8080` | No |

## Local Development

### Using Docker Compose

```bash
# Start PostgreSQL and app
docker compose up -d

# View logs
docker compose logs -f app

# Stop services
docker compose down

# Stop and remove volumes
docker compose down -v
```

### Without Docker

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your local PostgreSQL credentials

3. Run the application:
   ```bash
   go run ./cmd/app
   ```

## Nixpacks Configuration

The `nixpacks.toml` file configures the build:

```toml
[phases.setup]
nixPkgs = ["...", "postgresql"]

[phases.build]
cmd = "go build -ldflags '-s -w' -o ./bin/app ./cmd/app"

[start]
cmd = "./bin/app"

[variables]
CGO_ENABLED = "0"
GOOS = "linux"
```

### Custom Build Commands

Override in Dokploy environment variables if needed:

```bash
NIXPACKS_BUILD_CMD=go build -o main ./cmd/app
NIXPACKS_START_CMD=./main
```

## Database Backups

Configure automatic backups in Dokploy:

1. Go to **PostgreSQL Service** → **Backups**
2. Configure:
   - **Schedule**: `0 2 * * *` (daily at 2 AM)
   - **Prefix**: `graph-rag-backup`
   - **Database**: `graphrag`
   - **Retention**: 7 (keep last 7 backups)
3. Configure S3 destination in **Settings** → **S3 Destinations**

## Health Checks

The application exposes `/health` endpoint. Configure in Dokploy:

**Application** → **Advanced** → **Health Check**:

```yaml
Test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
Interval: 30s
Timeout: 10s
Retries: 3
Start Period: 40s
```

## Troubleshooting

### Connection Refused

- Ensure PostgreSQL is deployed and running (green status)
- Verify `POSTGRES_HOST` matches the database service name
- Check both services are in the same project

### Build Failures

- Check build logs in Dokploy UI
- Verify `go.mod` exists and dependencies are correct
- Test locally: `nixpacks build .`

### Database Migration Issues

- Connect to PostgreSQL container: `docker exec -it <container> psql -U graphrag -d graphrag`
- Check logs: `docker compose logs postgres`

## Resources

- [Dokploy Documentation](https://docs.dokploy.com)
- [Nixpacks Documentation](https://nixpacks.com/docs)
- [Nixpacks Go Provider](https://nixpacks.com/docs/providers/go)
