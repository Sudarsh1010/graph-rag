# Phase 4: pgvector Embeddings Implementation

## Objective
Add semantic search capabilities using pgvector to enable similarity-based queries on text fields.

## Tasks

### 1. Identify Text Fields for Embeddings
- [ ] product_descriptions.productDescription
- [ ] business_partners.businessPartnerFullName
- [ ] business_partner_addresses.{cityName, country, region, street}
- [ ] plants.plantName
- [ ] products.productGroup, productType
- [ ] Any other descriptive text fields identified in the data

### 2. Add Vector Columns
For each table with text to embed:
- [ ] Alter table to add `embedding vector(384)` (or 768 depending on model choice)
- [ ] Example: `ALTER TABLE product_descriptions ADD COLUMN embedding vector(384);`

### 3. Select and Set Up Embedding Model
- [ ] Choose a sentence-transformers model (e.g., all-MiniLM-L6-v2 for 384 dimensions)
- [ ] Alternatively, use OpenAI or other API if preferred and budget allows
- [ ] Implement embedding generation service in Go
- [ ] Handle batching for efficiency

### 4. Generate and Store Embeddings
- [ ] For each text field, generate embedding and store in the vector column
- [ ] Example: Update product_descriptions set embedding = <vector> where product = ? and language = ?
- [ ] Handle null/empty text (store zero vector or skip?)
- [ ] Consider updating embeddings when source text changes (if implementing incremental updates)

### 5. Create Vector Indexes
- [ ] Create ivfflat or hnsw index on each vector column for cosine similarity
- [ ] Example: `CREATE INDEX ON product_descriptions USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);`
- [ ] Monitor and adjust `lists` parameter based on data size

### 6. Implement Similarity Search Functions
- [ ] Create SQL functions or API endpoints for:
  - Find similar products by description
  - Find similar customers by name/address
  - Find similar plants by name
  - Cross-entity similarity (if useful)
- [ ] Example query: 
  ```sql
  SELECT product, language, productDescription, 
         1 - (embedding <=> <query_vector>) AS similarity
  FROM product_descriptions
  WHERE embedding IS NOT NULL
  ORDER BY embedding <=> <query_vector>
  LIMIT 10;
  ```

## Implementation Details

### Technology Stack
- **Embedding Generation**: Go program using a sentence-transformers binding or Python subprocess
- **Alternative**: Call external API (if using OpenAI, Cohere, etc.)
- **Vector Storage**: PostgreSQL pgvector extension
- **Indexing**: pgvector's ivfflat or hnsw

### File Structure
```
internal/
  embedding/
    generator.go     # Interface for embedding generation
    local.go         # Local sentence-transformers implementation
    # api.go         # If using external API
    service.go       # Service to generate embeddings for entities
    updater.go       # Job to update embeddings (batch or trigger-based)
```

### Atomic Commit Strategy
- **feat(embedding): add vector columns to identified tables**
- **feat(embedding): implement embedding generation service**
- **feat(embedding): add batch embedding update jobs**
- **feat(embedding): create vector indexes**
- **feat(embedding): implement similarity search functions in database or API**

## Dependencies
- Completion of Phase 2: Relational Schema Setup (tables must exist)
- Data loaded from Phase 1 (we need the text fields populated)
- Choice of embedding model and dimension

## Verification Criteria
After embedding implementation:
- [ ] Vector columns added to specified tables
- [ ] Embedding generation service runs without errors
- [ ] A sample of records have non-zero vector embeddings
- [ ] Vector indexes created successfully
- [ ] Similarity search queries return results and are ordered by similarity
- [ ] Search results are relevant (manual verification on sample queries)
- [ ] Performance: index search returns in acceptable time (<100ms for dataset size)

## Estimated Effort
- **Model selection and setup**: 2-4 hours
- **Embedding service**: 4-6 hours
- **Batch update jobs**: 2-4 hours
- **Indexing**: 2-4 hours
- **Search function implementation**: 2-4 hours
- **Total**: 2-3 days

## Next Steps
Upon completion, proceed to Phase 5: API Development to create REST endpoints for accessing the data, graph, and vector search.