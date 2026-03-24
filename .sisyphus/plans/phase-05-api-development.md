# Phase 5: API Development

## Objective
Create a Go-based REST API to expose the SAP data, graph traversals, and vector search capabilities.

## Tasks

### 1. Set Up API Framework
- [ ] Choose a web framework (e.g., Gin, Echo, or net/http if keeping it simple)
- [ ] Set up routing, middleware (logging, CORS, timeout)
- [ ] Configure database connection pooling
- [ ] Set up error handling and response formatting

### 2. Implement CRUD Endpoints for SAP Entities
For each of the 19 entities (or logical groupings):
- [ ] GET /api/v1/{entity} - List with pagination, filtering, sorting
- [ ] GET /api/v1/{entity}/{id} - Get single record
- [ ] POST /api/v1/{entity} - Create new record
- [ ] PUT /api/v1/{entity}/{id} - Update record
- [ ] DELETE /api/v1/{entity}/{id} - Delete record (if applicable, many may be read-only)
- [ ] Implement proper validation, authorization (if needed), and status codes

### 3. Implement Graph Query Endpoints
- [ ] POST /api/v1/graph/traverse - Execute a Cypher traversal query
  - Example body: { "query": "MATCH (n) RETURN n LIMIT 10", "params": {} }
- [ ] POST /api/v1/graph/path - Find shortest path between two nodes
  - Example body: { "start": { "label": "SalesOrder", "property": { "salesOrder": "740506" } }, "end": { "label": "Payment", "property": { "accountingDocument": "..." } }, "relationships": ["DELIVERS", "BILLS", "PAYS_FOR"] }
- [ ] GET /api/v1/graph/vertices/{label} - Get vertices by label with optional filtering
- [ ] GET /api/v1/graph/edges/{label} - Get edges by label with optional filtering
- [ ] Implement proper sanitization to prevent Cypher injection

### 4. Implement Vector Search Endpoints
- [ ] GET /api/v1/search/products?query=<text>&limit=10 - Search products by description similarity
- [ ] GET /api/v1/search/customers?query=<text>&limit=10 - Search customers by name/address similarity
- [ ] GET /api/v1/search/plants?query=<text>&limit=10 - Search plants by name similarity
- [ ] POST /api/v1/search/vector - Generic vector search: provide vector and search options
- [ ] Implement query parameter validation and sanitization

### 5. Implement Batch and Export Endpoints
- [ ] GET /api/v1/export/{entity}?format=csv - Export entity data as CSV
- [ ] POST /api/v1/batch/insert - Insert multiple records in a single transaction
- [ ] POST /api/v1/batch/update - Update multiple records

### 6. Add Health Checks and Metrics
- [ ] GET /health - Liveness and readiness checks
- [ ] GET /metrics - Prometheus metrics (if desired)
- [ ] Monitor database connection pool, query latency, error rates

### 7. Implement Authentication and Authorization (if required)
- [ ] Depending on requirements, add API key, JWT, or other auth mechanism
- [ ] Role-based access control for different entities or operations

### 8. API Documentation
- [ ] Generate OpenAPI/Swagger documentation
- [ ] Use annotations or separate files to document endpoints
- [ ] Make documentation available at /swagger-ui.html or similar

## Implementation Details

### Technology Stack
- **Language**: Go (1.19+)
- **Web Framework**: To be decided (Gin is popular and performant)
- **Database**: database/sql with PostgreSQL driver (pgx)
- **Graph Queries**: Use PostgreSQL's AGE extension via SQL/Cypher
- **Vector Search**: Use pgvector via SQL
- **Dependency Injection**: Possibly use wire or manual DI
- **Logging**: Structured logging (zap, zerolog)
- **Testing**: Go's built-in testing package, testify for assertions

### File Structure
```
cmd/
  server/
    main.go                 # Entry point
internal/
  api/
    handler.go              # HTTP handler setup
    middleware.go           # Custom middleware
    v1/
      sales_order.go        # Handlers for sales order endpoints
      customer.go           # Handlers for customer endpoints
      # ... etc for each entity or logical group
      graph.go              # Graph query handlers
      search.go             # Vector search handlers
  service/
    sales_order_service.go  # Business logic for sales orders
    # ... etc
    graph_service.go        # Wrapper for AGE queries
    search_service.go       # Wrapper for vector search
  store/
    sales_order_store.go    # Database access layer
    # ... etc
config/
  config.go                 # Configuration loading
```

### Atomic Commit Strategy
- **feat(api): set up project structure and basic server**
- **feat(api): add CRUD endpoints for [entity group]** (e.g., sales order, billing, delivery)
- **feat(api): add graph query endpoints**
- **feat(api): add vector search endpoints**
- **feat(api): add batch and export endpoints**
- **feat(api): add health checks and metrics**
- **feat(api): add authentication and authorization (if needed)**
- **feat(api): add API documentation (Swagger/OpenAPI)**

## Dependencies
- Completion of Phases 1-4: Data must be loaded into PostgreSQL, AGE graph, and vector columns populated
- Database connection details from configuration
- Graph and vector search functions available in the database

## Verification Criteria
After API implementation:
- [ ] Server starts successfully and listens on configured port
- [ ] CRUD endpoints return correct status codes and data
- [ ] Graph traversal endpoints execute Cypher queries and return results
- [ ] Vector search endpoints return results ordered by similarity
- [ ] Error handling returns appropriate HTTP status codes (400, 404, 500)
- [ ] Validation works for input data (e.g., reject invalid salesOrder format)
- [ ] Concurrent requests handled properly (test with load)
- [ ] Health check endpoint returns 200 OK when healthy
- [ ] API documentation is accessible and accurate

## Estimated Effort
- **Project setup and middleware**: 4-6 hours
- **CRUD endpoints**: 8-12 hours (19 entities, but many can be similar)
- **Graph endpoints**: 4-6 hours
- **Vector search endpoints**: 4-6 hours
- **Batch/export and health checks**: 4-6 hours
- **Authentication (if needed)**: 2-4 hours
- **Documentation**: 2-4 hours
- **Total**: 2-3 days

## Next Steps
Upon completion, proceed to Phase 6: Frontend Integration to build the React client that consumes these APIs.