# Phase 1: Data Ingestion & Loading

## Objective
Load 19 SAP order-to-cash JSONL entity types directly into PostgreSQL via the existing Go Bun models, performing validation and type conversion during the load.

## Tasks (Parallelizable by Entity Type)

### 1. Create Entity Loaders
For each of the 19 entities, create a loader function that:
- [ ] sales_order_headers loader
- [ ] sales_order_items loader  
- [ ] sales_order_schedule_lines loader
- [ ] billing_document_headers loader
- [ ] billing_document_items loader
- [ ] billing_document_cancellations loader
- [ ] outbound_delivery_headers loader
- [ ] outbound_delivery_items loader
- [ ] payments_accounts_receivable loader
- [ ] journal_entry_items_accounts_receivable loader
- [ ] business_partners loader
- [ ] business_partner_addresses loader
- [ ] customer_company_assignments loader
- [ ] customer_sales_area_assignments loader
- [ ] products loader
- [ ] product_descriptions loader
- [ ] product_plants loader
- [ ] product_storage_locations loader
- [ ] plants loader

Each loader will:
- Ensure the corresponding table exists (using Bun's `CreateTableIfNotExists` or via migration)
- Read all JSONL files in the corresponding entity directory
- Parse each line as JSON
- Convert string numbers to `decimal.Decimal`
- Parse ISO 8601 timestamps to `time.Time`
- Trim whitespace and handle empty strings as null for nullable fields
- Validate data (ranges, required fields, etc.)
- Insert valid records using Bun (via `db.NewInsert().Model(&model).Exec(ctx)`)
- Log statistics and errors

### 2. Implement Shared Validation & Conversion Helpers
- [ ] Number string to decimal conversion (handle empty strings, invalid formats)
- [ ] ISO timestamp parsing (handle invalid formats)
- [ ] Validation functions per entity (e.g., amount >= 0, valid status codes)
- [ ] Error logging and reporting

## Implementation Details

### Technology Stack
- **Language**: Go (using existing `internal/model/*` structs and Bun ORM)
- **JSON Processing**: `encoding/json` (with custom unmarshaling or manual conversion after parsing into map)
- **Validation**: Custom validation functions
- **Error Handling**: Return detailed errors; continue loading other records
- **Logging**: Structured logging (using existing zap logger from `internal/di/logger.go`)

### File Structure
```
internal/
  loader/
    sales_order_headers.go
    sales_order_items.go
    # ... etc for all 19 entities
    converter.go      # Shared number/date conversion
    validator.go      # Shared validation functions
    coordinator.go    # Coordinates loading all entities (can run in parallel)
```

### Atomic Commit Strategy
- **One commit per entity loader** (19 commits total)
  - Example: `feat(loader): add sales_order_headers loader with validation and insertion`
- **Shared utilities commits**
  - `feat(loader): add conversion and validation helpers`
  - `feat(loader): add loader coordinator to run entity loaders in parallel`

## Dependencies
- Existing Go models in `internal/model/`
- Bun database connection (from `internal/di/database.go`)

## Verification Criteria
After loading:
- [ ] All JSONL records are processed (attempted)
- [ ] Valid records are inserted into the corresponding PostgreSQL tables
- [ ] Invalid records are logged (with reason) and skipped
- [ ] Data types in the database match the model definitions (decimal, timestamp, etc.)
- [ ] Referential integrity is maintained (foreign keys point to existing records where applicable)
- [ ] Load process completes without panics
- [ ] All tables exist (created if not present)

## Estimated Effort
- **Per entity loader**: 1-2 hours (simplified: ensure table exists, read, convert, validate, insert)
- **Total**: 1 day for all 19 entities (parallelizable)
- **Shared utilities**: 2-4 hours
- **Loader coordinator**: 2-4 hours

## Next Steps
Upon completion, proceed to Phase 2: Apache AGE Graph Construction to build the graph from the loaded relational data.