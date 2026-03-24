# Phase 0: Infrastructure (Docker + Extensions)

## Objective
Set up the PostgreSQL Docker image with Apache AGE and pgvector extensions pre-installed, and update docker-compose.yml to use it.

## Tasks

### 1. Create Custom PostgreSQL Dockerfile
- [ ] Create `docker/Dockerfile.postgres` based on `pgvector/pgvector:pg16`
- [ ] Install Apache AGE extension (compile from source or use available package)
- [ ] Configure `shared_preload_libraries = 'age'` in postgresql.conf
- [ ] Ensure both extensions can be created: `CREATE EXTENSION vector; CREATE EXTENSION age;`

### 2. Update docker-compose.yml
- [ ] Replace `postgres:16-alpine` with custom build
- [ ] Update volume mounts if needed
- [ ] Update healthcheck
- [ ] Ensure app container can connect

### 3. Verify Extensions Work
- [ ] Start docker-compose
- [ ] Connect to postgres and run `CREATE EXTENSION vector;`
- [ ] Connect to postgres and run `CREATE EXTENSION age;`
- [ ] Test basic AGE operation: `SELECT * FROM create_graph('test');`
- [ ] Test basic pgvector operation: `SELECT '[1,2,3]'::vector;`

## Implementation Details

### Dockerfile (`docker/Dockerfile.postgres`)
```dockerfile
FROM pgvector/pgvector:pg16

# Install Apache AGE build dependencies and compile from source
# Or use a pre-built package if available for PG16
# Configure shared_preload_libraries
```

### docker-compose.yml Changes
```yaml
services:
  postgres:
    build:
      context: ./docker
      dockerfile: Dockerfile.postgres
    container_name: graph-rag-postgres
    environment:
      POSTGRES_USER: graphrag
      POSTGRES_PASSWORD: graphrag_dev_password
      POSTGRES_DB: graphrag
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U graphrag -d graphrag"]
      interval: 5s
      timeout: 5s
      retries: 5
```

### Atomic Commit Strategy
- **feat(docker): add custom PostgreSQL image with AGE + pgvector extensions**
- **feat(docker): update docker-compose to use custom Postgres image**

## Dependencies
- None (this is the foundation)

## Verification Criteria
- [ ] `docker-compose build` succeeds
- [ ] `docker-compose up` starts PostgreSQL
- [ ] `CREATE EXTENSION vector;` succeeds
- [ ] `CREATE EXTENSION age;` succeeds
- [ ] `SELECT create_graph('test');` returns successfully
- [ ] `SELECT '[1,2,3]'::vector;` returns a vector
- [ ] App container can connect to PostgreSQL

## Estimated Effort
- **Dockerfile creation**: 2-4 hours (AGE compilation can be tricky)
- **docker-compose update**: 30 min
- **Verification**: 1 hour
- **Total**: ~4-6 hours

## Next Steps
Upon completion, Phase 1 (Data Ingestion) and Phase 2 (Schema) can begin in parallel.
