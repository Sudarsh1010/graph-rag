# Context Graph System for SAP Order-to-Cash Dataset

## Project Overview
Build a context graph system using PostgreSQL with pgvector and Apache AGE extensions to load, query, and analyze SAP order-to-cash JSONL dataset. The system will enable LLM-powered natural language queries over interconnected SAP data.

## Key Components
1. **Data Ingestion Pipeline** - Process 19 JSONL entity types into relational tables
2. **PostgreSQL Database** - Store normalized SAP data with proper constraints
3. **Apache AGE Graph** - Enable graph traversals of order-to-cash lifecycle
4. **pgvector Extension** - Provide semantic search capabilities
5. **Go Backend API** - REST endpoints for data access and graph queries
6. **React Frontend** - UI for exploring data, graph visualizations, and search

## Success Criteria
- All 19 entity types loaded with referential integrity
- Apache AGE supports full order-to-cash lifecycle traversal
- pgvector enables semantic search with >80% precision
- API responds to 95% of requests under 200ms
- Frontend provides intuitive exploration of tabular and graph views
- Test coverage >80% for critical paths

## Execution Approach
Follow the detailed phase plans in order, with opportunities for parallel execution where noted. Each phase includes specific tasks, atomic commit strategies, and verification criteria.

## Related Plans
- [Phase 1: Data Ingestion & Preprocessing](./phase-01-data-ingestion.md)
- [Phase 2: Relational Schema Setup](./phase-02-schema-setup.md)
- [Phase 3: Apache AGE Graph Construction](./phase-03-age-graph.md)
- [Phase 4: pgvector Embeddings Implementation](./phase-04-pgvector-embeddings.md)
- [Phase 5: API Development](./phase-05-api-development.md)
- [Phase 6: Frontend Integration](./phase-06-frontend-integration.md)
- [Phase 7: Testing Strategy](./phase-07-testing-strategy.md)