# Phase 2: Relational Schema Setup

## Objective
Design and create PostgreSQL tables for all 19 SAP entities with proper data types, primary keys, foreign keys, and constraints.

## Tasks

### 1. Define Table Schemas
- [ ] sales_order_headers table
- [ ] sales_order_items table
- [ ] sales_order_schedule_lines table
- [ ] billing_document_headers table
- [ ] billing_document_items table
- [ ] billing_document_cancellations table
- [ ] outbound_delivery_headers table
- [ ] outbound_delivery_items table
- [ ] payments_accounts_receivable table
- [ ] journal_entry_items_accounts_receivable table
- [ ] business_partners table
- [ ] business_partner_addresses table
- [ ] customer_company_assignments table
- [ ] customer_sales_area_assignments table
- [ ] products table
- [ ] product_descriptions table
- [ ] product_plants table
- [ ] product_storage_locations table
- [ ] plants table

### 2. Implement Primary Keys
- [ ] Identify and set composite primary keys where appropriate
- [ ] Example: sales_order_items (salesOrder, salesOrderItem)
- [ ] Example: sales_order_schedule_lines (salesOrder, salesOrderItem, scheduleLine)
- [ ] Example: product_descriptions (product, language)

### 3. Establish Foreign Key Relationships
- [ ] sales_order_items.salesOrder → sales_order_headers.salesOrder
- [ ] sales_order_schedule_lines.(salesOrder, salesOrderItem) → sales_order_items.(salesOrder, salesOrderItem)
- [ ] outbound_delivery_items.referenceSdDocument → sales_order_headers.salesOrder (or delivery header)
- [ ] billing_document_items.referenceSdDocument → sales_order_headers.salesOrder (or delivery header)
- [ ] outbound_delivery_items.plant → plants.plant
- [ ] outbound_delivery_items.storageLocation → product_storage_locations.(product, plant, storageLocation) [Note: complex]
- [ ] product_descriptions.product → products.product
- [ ] product_plants.(product, plant) → products.product and plants.plant
- [ ] product_storage_locations.(product, plant, storageLocation) → products.product and plants.plant
- [ ] business_partner_addresses.businessPartner → business_partners.businessPartner
- [ ] customer_company_assignments.customer → business_partners.businessPartner (where customer=true)
- [ ] customer_sales_area_assignments.customer → business_partners.businessPartner (where customer=true)
- [ ] payments_accounts_receivable.customer → business_partners.businessPartner
- [ ] journal_entry_items_accounts_receivable.customer → business_partners.businessPartner
- [ ] journal_entry_items_accounts_receivable.referenceDocument → billing_document_headers.billingDocument
- [ ] payments_accounts_receivable.accountingDocument → journal_entry_items_accounts_receivable.accountingDocument

### 4. Add Constraints and Indexes
- [ ] NOT NULL constraints on critical fields
- [ ] CHECK constraints for data validation (e.g., amounts >= 0, valid status codes)
- [ ] UNIQUE constraints where applicable (e.g., businessPartnerFullName if applicable)
- [ ] Indexes on foreign key columns for join performance
- [ ] Indexes on frequently queried fields (dates, amounts, status codes)
- [ ] Indexes on reference fields (referenceSdDocument, etc.)

### 5. Data Types and Precision
- [ ] Use appropriate types: VARCHAR for codes, DECIMAL for amounts, TIMESTAMP for dates
- [ ] Set proper lengths based on SAP specifications (if known) or observed data
- [ ] Use TEXT for long descriptions if needed

## Implementation Details

### Technology Stack
- **SQL**: Standard PostgreSQL DDL
- **Migration Tool**: Possibly using existing migration system in codebase (check internal/ for migration patterns)
- **Alternatively**: Create SQL files to be executed via migration tool

### File Structure
```
# If using migrations:
migrations/
  001_create_sales_order_headers.sql
  002_create_sales_order_items.sql
  # ... etc
# Or single schema file:
database/
  schema.sql
```

### Atomic Commit Strategy
- **One commit per logical group** of related tables
  - Example: `feat(db): add sales order schema (headers, items, schedule lines)`
  - Example: `feat(db): add billing schema (headers, items, cancellations)`
  - Example: `feat(db): add delivery schema (headers, items)`
  - Example: `feat(db): add financial schema (payments, journal entries)`
  - Example: `feat(db): add master data schema (business partners, products, plants)`
- **Separate commit for indexes**: `feat(db): add performance indexes on foreign keys and query fields`
- **Separate commit for constraints**: `feat(db): add data validation constraints`

## Dependencies
- Completion of Phase 1: Data Ingestion (to know exact data types and validation needs)
- Existing database package in internal/ if available

## Verification Criteria
After schema setup:
- [ ] All 19 tables created successfully
- [ ] Primary keys correctly defined
- [ ] Foreign key constraints created and functional
- [ ] Check constraints prevent invalid data
- [ ] Indexes created on expected columns
- [ ] Data types match the cleaned data from Phase 1
- [ ] Can insert sample data without constraint violations
- [ ] Can perform basic joins using foreign keys

## Estimated Effort
- **Schema design**: 4-6 hours (analyzing relationships and data types)
- **Implementation**: 2-4 hours (writing DDL statements)
- **Total**: 1 day

## Next Steps
Upon completion, proceed to Phase 3: Apache AGE Graph Construction to define graph schema and load data into the graph.