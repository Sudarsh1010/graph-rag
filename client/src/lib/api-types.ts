// Types matching the actual snake_case column names returned by the Go API.
// The API uses scanRows which returns raw DB column names.

export interface SalesOrderHeader {
  sales_order: string
  sales_order_type: string
  sales_organization: string
  distribution_channel: string
  organization_division: string
  sales_group: string | null
  sales_office: string | null
  sold_to_party: string
  creation_date: string
  created_by_user: string
  last_change_date_time: string
  total_net_amount: string
  overall_delivery_status: string
  overall_ord_reltd_billg_status: string | null
  overall_sd_doc_reference_status: string | null
  total_credit_check_status: string | null
  transaction_currency: string
  pricing_date: string
  requested_delivery_date: string
  header_billing_block_reason: string | null
  delivery_block_reason: string | null
  incoterms_classification: string
  incoterms_location1: string | null
  customer_payment_terms: string
  id: number
  created_at: string
  updated_at: string
}

export interface SalesOrderItem {
  sales_order: string
  sales_order_item: string
  sales_order_item_category: string
  material: string
  requested_quantity: string
  requested_quantity_unit: string
  transaction_currency: string
  net_amount: string
  material_group: string
  production_plant: string
  storage_location: string
  id: number
  created_at: string
  updated_at: string
}

export interface Product {
  product: string
  product_type: string
  cross_plant_status: string
  creation_date: string
  created_by_user: string
  last_change_date: string
  product_old_id: string
  gross_weight: string
  weight_unit: string
  net_weight: string
  product_group: string
  base_unit: string
  division: string
  industry_sector: string
  id: number
  created_at: string
  updated_at: string
}

export interface ProductDescription {
  product: string
  language: string
  product_description: string
  id: number
  created_at: string
  updated_at: string
}

export interface Customer {
  business_partner: string
  customer: string
  business_partner_category: string
  business_partner_full_name: string
  business_partner_grouping: string
  business_partner_name: string
  correspondence_language: string | null
  created_by_user: string
  creation_date: string
  creation_time: { hour: number; minute: number; second: number }
  first_name: string | null
  form_of_address: string
  industry: string | null
  last_change_date: string
  last_name: string | null
  organization_bp_name1: string
  organization_bp_name2: string | null
  business_partner_is_blocked: boolean
  is_marked_for_archiving: boolean
  id: number
  created_at: string
  updated_at: string
}

export interface Plant {
  plant: string
  plant_name: string
  valuation_area: string
  plant_customer: string
  plant_supplier: string
  factory_calendar: string
  default_purchasing_organization: string
  sales_organization: string
  address_id: string
  plant_category: string
  distribution_channel: string
  division: string
  language: string
  is_marked_for_archiving: boolean
  id: number
  created_at: string
  updated_at: string
}

export interface BillingDocumentHeader {
  billing_document: string
  billing_document_type: string
  company_code: string
  fiscal_year: string
  creation_date: string
  creation_time: { hour: number; minute: number; second: number }
  last_change_date_time: string | null
  billing_document_date: string
  billing_document_is_cancelled: boolean | null
  cancelled_billing_document: string | null
  total_net_amount: string
  transaction_currency: string
  accounting_document: string | null
  sold_to_party: string
  id: number
  created_at: string
  updated_at: string
}

export interface Delivery {
  actual_goods_movement_date: string | null
  actual_goods_movement_time: { hour: number; minute: number; second: number }
  creation_date: string
  creation_time: { hour: number; minute: number; second: number }
  delivery_block_reason: string
  delivery_document: string
  hdr_general_incompletion_status: string
  header_billing_block_reason: string
  last_change_date: string | null
  overall_goods_movement_status: string
  overall_picking_status: string
  overall_proof_of_delivery_status: string
  shipping_point: string
  id: number
  created_at: string
  updated_at: string
}

export interface Payment {
  company_code: string
  fiscal_year: number
  accounting_document: string
  accounting_document_item: string
  clearing_date: string | null
  clearing_accounting_document: string
  clearing_doc_fiscal_year: number | null
  amount_in_transaction_currency: string
  transaction_currency: string
  amount_in_company_code_currency: string
  company_code_currency: string
  customer: string
  invoice_reference: string | null
  sales_document: string | null
  sales_document_item: string | null
  posting_date: string | null
  document_date: string | null
  assignment_reference: string | null
  gl_account: string
  financial_account_type: string
  profit_center: string
  cost_center: string | null
}

export interface JournalEntry {
  company_code: string
  fiscal_year: number
  accounting_document: string
  gl_account: string
  reference_document: string
  cost_center: string | null
  profit_center: string
  transaction_currency: string
  amount_in_transaction_currency: string
  company_code_currency: string
  amount_in_company_code_currency: string
  posting_date: string | null
  document_date: string | null
  accounting_document_type: string
  accounting_document_item: string
  assignment_reference: string | null
  last_change_date_time: string
  customer: string
  financial_account_type: string
  clearing_date: string | null
  clearing_accounting_document: string
  clearing_doc_fiscal_year: number | null
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
}

export interface SearchResponse {
  query: string
  results: Record<string, Record<string, unknown>[]>
}
