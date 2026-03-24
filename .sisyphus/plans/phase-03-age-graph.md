# Phase 3: Apache AGE Graph Construction

## Objective
Load the cleaned SAP data into Apache AGE to enable graph traversals and relationship-based queries.

## Tasks

### 1. Initialize Apache AGE Graph
- [ ] Ensure AGE extension is enabled in PostgreSQL (should be via docker-compose)
- [ ] Create a new graph: `SELECT * FROM create_graph('sap_ocg');`

### 2. Define Vertex Labels (One per Entity)
- [ ] Customer (from business_partners where customer=true)
- [ ] Vendor (from business_partners where vendor=true) [if applicable]
- [ ] Product
- [ ] Plant
- [ ] StorageLocation (composite: product+plant+storageLocation)
- [ ] SalesOrder
- [ ] SalesOrderItem
- [ ] SalesOrderScheduleLine
- [ ] DeliveryHeader
- [ ] DeliveryItem
- [ ] BillingHeader
- [ ] BillingItem
- [ ] Payment
- [ ] JournalEntry
- [ ] CompanyCode
- [ ] SalesOrganization
- [ ] DistributionChannel
   *(Note: Some of these may be better as properties on vertices rather than separate vertices)*

### 3. Define Edge Labels (Relationships)
- [ ] CONTAINS: SalesOrder -> SalesOrderItem
- [ ] CONTAINS: SalesOrderItem -> SalesOrderScheduleLine
- [ ] REFERENCES: SalesOrderItem -> Product (via material)
- [ ] PRODUCED_AT: SalesOrderItem -> Plant (via productionPlant)
- [ ] SHIPPED_FROM: DeliveryItem -> Plant
- [ ] STORED_AT: DeliveryItem -> StorageLocation (via product, plant, storageLocation)
- [ ] DELIVERS: SalesOrder -> DeliveryHeader (via referenceSdDocument in delivery header? Actually delivery items reference SD document)
- [ ] BILLS: DeliveryHeader -> BillingHeader (via referenceSdDocument in billing header)
- [ ] ITEMIZES: BillingHeader -> BillingItem
- [ ] REFERS_TO: BillingItem -> SalesOrderItem (via referenceSdDocument in billing item)
- [ ] PAYS_FOR: Payment -> BillingHeader (via accountingDocument? Actually payment clears invoice)
- [ ] CLEARS: Payment -> JournalEntry (via accountingDocument)
- [ ] HAS_ADDRESS: BusinessPartner -> Address (if we model addresses separately)
- [ ] BELONGS_TO_COMPANY: Customer -> CompanyCode
- [ ] ASSIGNED_TO_SALES_ORG: Customer -> SalesOrganization
- [ ] ASSIGNED_TO_DIST_CHANNEL: Customer -> DistributionChannel
- [ ] HAS_PLANT_ASSIGNMENT: Product -> Plant (via product_plants)
- [ ] HAS_STORAGE_LOCATION: Product -> StorageLocation (via product_storage_locations)
- [ ] IS_LOCATED_AT: Plant -> Address (if we have address for plants)
- [ ] BELONGS_TO_BUSINESS_PARTNER: Address -> BusinessPartner

### 4. Load Vertices
For each vertex label:
- [ ] Insert vertices from corresponding table
- [ ] Use appropriate properties (e.g., Customer vertex: businessPartner, customer, businessPartnerFullName, etc.)
- [ ] Handle composite keys by creating a surrogate key or using the composite as id in AGE (AGE requires a single id, so we may need to create a surrogate or concatenate)

### 5. Load Edges
For each edge label:
- [ ] Insert edges by joining the relevant tables to get the source and target vertex IDs
- [ ] Example for CONTAINS (SalesOrder -> SalesOrderItem):
  ```sql
  INSERT INTO sap_ocg.edge (eid, start_vid, end_vid, label, properties)
  SELECT 
    nextval('sap_ocg.edge_id_seq') as eid,
    so.vid as start_vid,
    soi.vid as end_vid,
    'CONTAINS' as label,
    '{}'::agtype as properties
  FROM sales_order_header_vertices so
  JOIN sales_order_item_vertices soi ON so.salesOrder = soi.salesOrder
  ```
  (We'll need to create temporary views or tables that map our relational IDs to AGE vertex IDs)

### 6. Create Graph Indexes
- [ ] Indexes on vertex properties for fast lookups (e.g., Customer.businessPartner, Product.product)
- [ ] Indexes on edge labels for fast traversal filtering

## Implementation Details

### Technology Stack
- **SQL/PG**: Apache AGE cypher and SQL procedures
- **Approach**: Use SQL to insert into AGE's graph tables (agn_catalog.edge, agn_catalog.vertex, etc.) or use cypher loading
- **Alternative**: Use `cypher_load` CSV files if we export the graph data

### File Structure
```
internal/
  graph/
    initializer.go   # Sets up the graph, labels
    loader.go        # Loads vertices and edges
    queries.go       # Common graph traversal functions
```

### Atomic Commit Strategy
- **feat(graph): initialize AGE graph and create vertex/edge labels**
- **feat(graph): load vertices for [entity group]** (can split by logical group)
- **feat(graph): load edges for [relationship group]**
- **feat(graph): add performance indexes on graph properties**

## Dependencies
- Completion of Phase 2: Relational Schema Setup (we need the tables populated with data from Phase 1)
- Data from Phase 1 must be loaded into the PostgreSQL tables

## Verification Criteria
After graph construction:
- [ ] Graph created successfully: `SELECT * FROM ag_catalog.graph;`
- [ ] All vertex and edge labels created
- [ ] Vertex counts match source table counts (after any filtering)
- [ ] Edge counts match expected relationships
- [ ] Sample traversals work:
  - Find all items for a sales order: `MATCH (so:SalesOrder)-[:CONTAINS]->(item:SalesOrderItem) WHERE so.salesOrder = '740506' RETURN item;`
  - Trace order-to-cash: `MATCH path = (so:SalesOrder)-[:CONTAINS*]->(:SalesOrderItem)-[:REFERENCES*]->(:Product)-[:HAS_PLANT_ASSIGNMENT]->(p:Plant) RETURN so, p LIMIT 10;`
  - Find delivery for a sales order: `MATCH (so:SalesOrder)-[:DELIVERS]->(dh:DeliveryHeader) WHERE so.salesOrder = '740506' RETURN dh;`
  - Find payment for a billing document: `MATCH (bh:BillingHeader)<-[:PAYS_FOR]-(p:Payment) WHERE bh.billingDocument = 'XXXXXX' RETURN p;`
- [ ] No orphan vertices or edges (unless intentional)

## Estimated Effort
- **Graph setup and labeling**: 2-4 hours
- **Vertex loading**: 4-6 hours (19 entities, but some are small)
- **Edge loading**: 4-6 hours (multiple relationships, requires joins)
- **Indexing and testing**: 2-4 hours
- **Total**: 2-3 days

## Next Steps
Upon completion, proceed to Phase 4: pgvector Embeddings Implementation to add semantic search capabilities.