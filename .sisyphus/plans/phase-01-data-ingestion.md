# Phase 1: Data Ingestion & Preprocessing

## Objective
Process 19 SAP order-to-cash JSONL entity types into clean, validated data ready for relational loading.

## Tasks (Parallelizable by Entity Type)

### 1. Create JSONL Parsers
- [ ] sales_order_headers parser
- [ ] sales_order_items parser  
- [ ] sales_order_schedule_lines parser
- [ ] billing_document_headers parser
- [ ] billing_document_items parser
- [ ] billing_document_cancellations parser
- [ ] outbound_delivery_headers parser
- [ ] outbound_delivery_items parser
- [ ] payments_accounts_receivable parser
- [ ] journal_entry_items_accounts_receivable parser
- [ ] business_partners parser
- [ ] business_partner_addresses parser
- [ ] customer_company_assignments parser
- [ ] customer_sales_area_assignments parser
- [ ] products parser
- [ ] product_descriptions parser
- [ ] product_plants parser
- [ ] product_storage_locations parser
- [ ] plants parser

### 2. Implement Data Cleaning & Validation
For each parser:
- [ ] Handle null values appropriately per field
- [ ] Parse ISO 8601 timestamps to time.Time
- [ ] Convert decimal strings to proper numeric types
- [ ] Trim whitespace from string fields
- [ ] Validate reference IDs exist in parent entities (where applicable)
- [ ] Check data ranges (amounts >= 0, valid status codes, etc.)
- [ ] Log malformed records for review

### 3. Generate Intermediate Format
- [ ] Output cleaned data as JSONL or CSV per entity
- [ ] Alternatively, generate SQL INSERT statements
- [ ] Include metadata: source file, processing timestamp, record count

## Implementation Details

### Technology Stack
- **Language**: Go (leveraging existing cmd/ and internal/ packages)
- **JSON Processing**: encoding/json or jsoniter
- **Validation**: Custom validation functions per entity
- **Error Handling**: Return detailed errors for malformed records
- **Logging**: Structured logging for audit trail

### File Structure
```
internal/
  etl/
    parsers/
      sales_order_headers.go
      sales_order_items.go
      # ... etc for all 19 entities
    cleaner.go
    validator.go
    loader.go
```

### Atomic Commit Strategy
- **One commit per entity parser** (19 commits total)
  - Example: `feat(etl): add sales_order_headers parser with validation`
- **Shared utilities commits**
  - `feat(etl): add data cleaning and validation helpers`
  - `feat(etl): add ETL coordinator and batch processor`

## Dependencies
- None beyond standard library (keep it simple for reliability)
- May use existing database/internal packages if available

## Verification Criteria
Each parser must:
- [ ] Process 100% of input records without panics
- [ ] Produce correct output format (JSONL/CSV/SQL)
- [ ] Log statistics: total processed, valid, invalid records
- [ ] Handle edge cases: empty files, malformed JSON, missing fields
- [ ] Pass unit tests with sample data

## Estimated Effort
- **Per parser**: 2-4 hours (depending on complexity)
- **Total**: 1-2 days for all 19 entities (parallelizable)
- **Shared utilities**: 4-6 hours

## Next Steps
Upon completion, proceed to Phase 2: Relational Schema Setup to create database tables matching the parsed data structure.