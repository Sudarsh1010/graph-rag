# Phase 8: NL Query System (LLM-Powered Natural Language)

## Objective
Enable users to ask natural language questions about SAP data and receive answers with relevant data (tables, graph traversals, summaries).

## Current State
- ❌ No LLM integration code exists
- ❌ No NL query endpoint exists
- ✅ Frontend plan includes chat interface (Phase 6)
- ✅ All data access via API (Phase 5)

## Design

### Approach: LLM-as-Query-Translator
The LLM translates natural language questions into structured API calls or Cypher queries. This is NOT a RAG system — the LLM doesn't need to read all the data. It needs to understand the schema and generate appropriate queries.

### Architecture
```
User Question
    ↓
[LLM] → Understand intent → Generate query plan
    ↓
[Query Executor] → Execute API calls / Cypher / SQL
    ↓
[LLM] → Format results into natural language answer
    ↓
Response (answer + data + visualization hints)
```

### Query Types the LLM Must Handle
1. **Entity lookup**: "Show me sales order 740506" → `GET /api/v1/sales-orders/740506`
2. **Filtered search**: "Find all orders over $10,000" → `GET /api/v1/sales-orders?minAmount=10000`
3. **Graph traversal**: "Trace the order-to-cash lifecycle for order 740506" → `GET /api/v1/graph/order/740506/lifecycle`
4. **Semantic search**: "Find products similar to laptop computers" → `GET /api/v1/search/products?query=laptop computers`
5. **Aggregation**: "What's the total revenue this month?" → Custom aggregation query
6. **Multi-step**: "Who are the top 5 customers by revenue, and what did they order?" → Multiple API calls

## Tasks

### 1. Design Schema Context for LLM
- [ ] Create a schema description document that the LLM can reference
- [ ] Include: entity names, field names, relationships, available API endpoints
- [ ] Format as structured system prompt

### 2. Implement NL Query Service — `internal/nlquery/`
- [ ] `service.go` — Main NL query orchestrator
- [ ] `translator.go` — LLM → structured query translation
- [ ] `executor.go` — Execute translated queries against API/graph
- [ ] `formatter.go` — Format results into natural language answers
- [ ] `schema.go` — Schema context for LLM prompts

### 3. Implement Query Translation
- [ ] Define structured query format (JSON schema)
- [ ] LLM generates query plan: `{type: "api_call" | "cypher" | "aggregate", params: {...}}`
- [ ] Validate generated queries before execution (prevent injection)
- [ ] Handle ambiguous queries by asking clarifying questions

### 4. Implement Query Execution
- [ ] Execute API calls using internal HTTP client (same as external API)
- [ ] Execute Cypher queries via AGE
- [ ] Execute aggregation queries via SQL
- [ ] Combine results from multiple queries

### 5. Implement Result Formatting
- [ ] LLM formats structured results into natural language
- [ ] Include data tables in response (for frontend rendering)
- [ ] Include visualization hints (chart type, graph traversal path)
- [ ] Cite sources (which entities/documents the answer comes from)

### 6. Add NL Query API Endpoint — `internal/api/v1/nlquery.go`
- [ ] `POST /api/v1/nlquery` — Ask a question
  ```json
  {
    "question": "What is the order-to-cash cycle time for our top 5 customers?",
    "conversation_id": "optional-for-context"
  }
  ```
- [ ] Response includes: answer text, data results, visualization hints, sources

### 7. Conversation Context (Optional Enhancement)
- [ ] Store conversation history for follow-up questions
- [ ] Support "show me more details about that last order" type queries
- [ ] In-memory or DB-backed conversation store

## Implementation Details

### File Structure
```
internal/
  nlquery/
    service.go         # NLQueryService — main orchestrator
    translator.go      # LLM prompt construction + response parsing
    executor.go        # Query execution (API, Cypher, SQL)
    formatter.go       # Result → natural language formatting
    schema.go          # Schema context for system prompt
    types.go           # Structured query/response types
```

### LLM Provider
- [ ] Use OpenAI API (GPT-4o-mini for cost efficiency)
- [ ] Or local LLM (Ollama) for offline/privacy
- [ ] Abstract behind interface for provider swapping

### Schema Context (System Prompt)
```
You are a SAP data analyst assistant. The system has the following entities and relationships:

ENTITIES:
- SalesOrder: salesOrder (PK), soldToParty, totalNetAmount, creationDate, status
- SalesOrderItem: salesOrder + salesOrderItem (PK), material, netAmount, plant
- BillingDocument: billingDocument (PK), totalNetAmount, creationDate, soldToParty
- Product: product (PK), productGroup, baseUnit, division
- Customer: businessPartner (PK), fullName, category
- Plant: plant (PK), name
...

RELATIONSHIPS:
- SalesOrder → SalesOrderItem (CONTAINS_ITEM)
- SalesOrderItem → Product (REFERENCES_PRODUCT)
- SalesOrder → Delivery (DELIVERS)
- Delivery → BillingDocument (BILLS)
- BillingDocument → Payment (PAYS)
...

AVAILABLE QUERIES:
1. API calls: GET /api/v1/{entity}?page=&limit=&filter=
2. Graph traversal: MATCH (n) -[:EDGE*]→ (m) WHERE ...
3. Semantic search: vector similarity on product_descriptions, business_partners, plants

Generate a structured JSON query plan to answer the user's question.
```

### Structured Query Format
```typescript
{
  "steps": [
    {
      "type": "api_call",
      "method": "GET",
      "path": "/api/v1/sales-orders",
      "params": {"sort": "totalNetAmount", "order": "desc", "limit": 5}
    },
    {
      "type": "aggregate",
      "sql": "SELECT sold_to_party, SUM(total_net_amount) FROM sales_order_headers GROUP BY sold_to_party ORDER BY sum DESC LIMIT 5"
    }
  ],
  "answer_format": "table_with_summary"
}
```

### Response Format
```json
{
  "answer": "The top 5 customers by total order value are: Customer A ($150K), Customer B ($120K), ...",
  "data": [
    {"customer": "Customer A", "totalOrders": 15, "totalRevenue": 150000},
    ...
  ],
  "visualization": {"type": "bar_chart", "x": "customer", "y": "totalRevenue"},
  "sources": ["sales_order_headers"]
}
```

### Atomic Commit Strategy
- `feat(nlquery): add schema context and LLM service interface`
- `feat(nlquery): implement query translation with OpenAI`
- `feat(nlquery): implement query executor (API + Cypher + SQL)`
- `feat(nlquery): implement result formatter`
- `feat(nlquery): add POST /api/v1/nlquery endpoint`
- `feat(frontend): add NL query chat interface to frontend`

## Dependencies
- Phases 1-6 complete (full system operational)
- OpenAI API key (or local LLM setup)
- Phase 6 frontend chat interface

## Verification Criteria
- [ ] Simple lookup queries work: "Show order 740506"
- [ ] Filter queries work: "Orders over $10,000"
- [ ] Graph traversals work: "Trace order 740506 to payment"
- [ ] Semantic search works: "Products like laptops"
- [ ] Aggregation works: "Total revenue by customer"
- [ ] Multi-step queries work: "Top customers and their orders"
- [ ] LLM responses are accurate and cite data sources
- [ ] Invalid/dangerous queries are rejected gracefully
- [ ] Response time < 5 seconds for most queries
- [ ] Conversation context works for follow-up questions

## Estimated Effort
- **Schema context + prompt design**: 4-6 hours
- **LLM service + translation**: 6-8 hours
- **Query executor**: 4-6 hours
- **Result formatter**: 4-6 hours
- **API endpoint**: 2-3 hours
- **Frontend chat UI**: 4-6 hours
- **Testing + iteration**: 4-6 hours
- **Total**: 4-5 days

## Next Steps
Upon completion, the system fully supports natural language queries over SAP data. Consider:
- [ ] Performance optimization (caching, streaming)
- [ ] Multi-user conversation management
- [ ] Fine-tuned LLM for SAP domain
- [ ] Export results (PDF, CSV)
