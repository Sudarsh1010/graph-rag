package graph

type vertexDef struct {
	label   string
	table   string
	columns []string
}

var vertexDefs = []vertexDef{
	{label: "Customer", table: "business_partners", columns: []string{"business_partner", "business_partner_name", "business_partner_category"}},
	{label: "Product", table: "products", columns: []string{"product", "product_type", "product_group", "base_unit", "division"}},
	{label: "Plant", table: "plants", columns: []string{"plant", "plant_name", "plant_category"}},
	{label: "SalesOrder", table: "sales_order_headers", columns: []string{"sales_order", "sales_order_type", "sold_to_party", "total_net_amount", "transaction_currency", "overall_delivery_status"}},
	{label: "SalesOrderItem", table: "sales_order_items", columns: []string{"sales_order", "sales_order_item", "material", "production_plant", "net_amount"}},
	{label: "Delivery", table: "outbound_delivery_headers", columns: []string{"delivery_document", "overall_goods_movement_status"}},
	{label: "DeliveryItem", table: "outbound_delivery_items", columns: []string{"delivery_document", "delivery_document_item", "reference_sd_document", "plant"}},
	{label: "BillingDocument", table: "billing_document_headers", columns: []string{"billing_document", "billing_document_type", "sold_to_party", "total_net_amount", "transaction_currency"}},
	{label: "BillingItem", table: "billing_document_items", columns: []string{"billing_document", "billing_document_item", "material", "net_amount"}},
	{label: "Payment", table: "payments_accounts_receivable", columns: []string{"fiscal_year", "accounting_document", "customer", "amount_in_transaction_currency"}},
	{label: "JournalEntry", table: "journal_entry_items_accounts_receivable", columns: []string{"fiscal_year", "accounting_document", "customer", "gl_account", "amount_in_transaction_currency"}},
}

type edgeDef struct {
	label string
	from  string
	to    string
	match string
}

var edgeDefs = []edgeDef{
	{label: "CONTAINS_ITEM", from: "SalesOrder", to: "SalesOrderItem", match: "a.sales_order = b.sales_order"},
	{label: "REFERENCES_PRODUCT", from: "SalesOrderItem", to: "Product", match: "a.material = b.product"},
	{label: "PRODUCED_AT", from: "SalesOrderItem", to: "Plant", match: "a.production_plant = b.plant"},
	{label: "CONTAINS_DELIVERY_ITEM", from: "Delivery", to: "DeliveryItem", match: "a.delivery_document = b.delivery_document"},
	{label: "DELIVERY_FOR_ORDER", from: "DeliveryItem", to: "SalesOrder", match: "a.reference_sd_document = b.sales_order"},
	{label: "DELIVERY_ITEM_PLANT", from: "DeliveryItem", to: "Plant", match: "a.plant = b.plant"},
	{label: "CONTAINS_BILLING_ITEM", from: "BillingDocument", to: "BillingItem", match: "a.billing_document = b.billing_document"},
	{label: "BILLING_ITEM_PRODUCT", from: "BillingItem", to: "Product", match: "a.material = b.product"},
	{label: "BILLS_CUSTOMER", from: "BillingDocument", to: "Customer", match: "a.sold_to_party = b.business_partner"},
	{label: "SOLD_TO", from: "SalesOrder", to: "Customer", match: "a.sold_to_party = b.business_partner"},
	{label: "PAYS_CUSTOMER", from: "Payment", to: "Customer", match: "a.customer = b.business_partner"},
}
