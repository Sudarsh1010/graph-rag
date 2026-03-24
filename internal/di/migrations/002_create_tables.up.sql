CREATE TABLE IF NOT EXISTS plants (
    plant VARCHAR(4) PRIMARY KEY,
    plant_name TEXT NOT NULL DEFAULT '',
    valuation_area TEXT NOT NULL DEFAULT '',
    plant_customer TEXT NOT NULL DEFAULT '',
    plant_supplier TEXT NOT NULL DEFAULT '',
    factory_calendar TEXT NOT NULL DEFAULT '',
    default_purchasing_organization TEXT NOT NULL DEFAULT '',
    sales_organization TEXT NOT NULL DEFAULT '',
    address_id TEXT NOT NULL DEFAULT '',
    plant_category TEXT NOT NULL DEFAULT '',
    distribution_channel TEXT NOT NULL DEFAULT '',
    division TEXT NOT NULL DEFAULT '',
    language TEXT NOT NULL DEFAULT '',
    is_marked_for_archiving BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS products (
    product VARCHAR(40) PRIMARY KEY,
    product_type VARCHAR(10) NOT NULL DEFAULT '',
    cross_plant_status VARCHAR(1) NOT NULL DEFAULT '',
    cross_plant_status_validity_date TIMESTAMPTZ,
    creation_date TIMESTAMPTZ NOT NULL,
    created_by_user VARCHAR(10) NOT NULL DEFAULT '',
    last_change_date TIMESTAMPTZ NOT NULL,
    last_change_date_time TIMESTAMPTZ NOT NULL,
    is_marked_for_deletion BOOLEAN,
    product_old_id VARCHAR(40) NOT NULL DEFAULT '',
    gross_weight DECIMAL(18,4) NOT NULL DEFAULT 0,
    weight_unit VARCHAR(3) NOT NULL DEFAULT '',
    net_weight DECIMAL(18,4) NOT NULL DEFAULT 0,
    product_group VARCHAR(10) NOT NULL DEFAULT '',
    base_unit VARCHAR(4) NOT NULL DEFAULT '',
    division VARCHAR(2) NOT NULL DEFAULT '',
    industry_sector VARCHAR(1) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS product_descriptions (
    product VARCHAR(40) NOT NULL,
    language VARCHAR(2) NOT NULL,
    product_description VARCHAR(300) NOT NULL DEFAULT '',
    PRIMARY KEY (product, language),
    CONSTRAINT fk_product_desc_product FOREIGN KEY (product) REFERENCES products(product)
);

CREATE TABLE IF NOT EXISTS product_plants (
    product VARCHAR(40) NOT NULL,
    plant VARCHAR(4) NOT NULL,
    country_of_origin VARCHAR(3) NOT NULL DEFAULT '',
    region_of_origin VARCHAR(4) NOT NULL DEFAULT '',
    production_invtry_managed_loc VARCHAR(40) NOT NULL DEFAULT '',
    availability_check_type VARCHAR(2) NOT NULL DEFAULT '',
    fiscal_year_variant VARCHAR(2) NOT NULL DEFAULT '',
    profit_center VARCHAR(10) NOT NULL DEFAULT '',
    mrp_type VARCHAR(2) NOT NULL DEFAULT '',
    PRIMARY KEY (product, plant),
    CONSTRAINT fk_product_plant_product FOREIGN KEY (product) REFERENCES products(product),
    CONSTRAINT fk_product_plant_plant FOREIGN KEY (plant) REFERENCES plants(plant)
);

CREATE TABLE IF NOT EXISTS product_storage_locations (
    product VARCHAR(40) NOT NULL,
    plant VARCHAR(4) NOT NULL,
    storage_location VARCHAR(10) NOT NULL,
    physical_inventory_block_ind VARCHAR(1) NOT NULL DEFAULT '',
    date_of_last_posted_cnt_un_rstrcd_stk TIMESTAMPTZ,
    PRIMARY KEY (product, plant, storage_location),
    CONSTRAINT fk_psl_product FOREIGN KEY (product) REFERENCES products(product),
    CONSTRAINT fk_psl_plant FOREIGN KEY (plant) REFERENCES plants(plant)
);

CREATE TABLE IF NOT EXISTS business_partners (
    business_partner VARCHAR(10) PRIMARY KEY,
    customer VARCHAR(10) NOT NULL DEFAULT '',
    business_partner_category VARCHAR(1) NOT NULL DEFAULT '',
    business_partner_full_name TEXT NOT NULL DEFAULT '',
    business_partner_grouping VARCHAR(4) NOT NULL DEFAULT '',
    business_partner_name TEXT NOT NULL DEFAULT '',
    correspondence_language VARCHAR(2),
    created_by_user VARCHAR(12) NOT NULL DEFAULT '',
    creation_date TIMESTAMPTZ NOT NULL,
    creation_time JSONB,
    first_name TEXT,
    form_of_address VARCHAR(4) NOT NULL DEFAULT '',
    industry TEXT,
    last_change_date TIMESTAMPTZ NOT NULL,
    last_name TEXT,
    organization_bp_name1 TEXT NOT NULL DEFAULT '',
    organization_bp_name2 TEXT,
    business_partner_is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    is_marked_for_archiving BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS business_partner_addresses (
    business_partner VARCHAR(10) NOT NULL,
    address_id VARCHAR(10) NOT NULL,
    validity_start_date TIMESTAMPTZ NOT NULL,
    validity_end_date TIMESTAMPTZ NOT NULL,
    address_uuid TEXT NOT NULL DEFAULT '',
    address_time_zone TEXT NOT NULL DEFAULT '',
    city_name TEXT,
    country TEXT NOT NULL DEFAULT '',
    po_box TEXT,
    po_box_deviating_city_name TEXT,
    po_box_deviating_country TEXT,
    po_box_deviating_region TEXT,
    po_box_is_without_number BOOLEAN NOT NULL DEFAULT FALSE,
    po_box_lobby_name TEXT,
    po_box_postal_code TEXT,
    postal_code TEXT,
    region TEXT,
    street_name TEXT,
    tax_jurisdiction TEXT,
    transport_zone TEXT,
    PRIMARY KEY (business_partner, address_id),
    CONSTRAINT fk_bpa_business_partner FOREIGN KEY (business_partner) REFERENCES business_partners(business_partner)
);

CREATE TABLE IF NOT EXISTS customer_company_assignments (
    customer VARCHAR(10) NOT NULL,
    company_code VARCHAR(4) NOT NULL,
    accounting_clerk TEXT,
    accounting_clerk_fax_number TEXT,
    accounting_clerk_internet_address TEXT,
    accounting_clerk_phone_number TEXT,
    alternative_payer_account TEXT,
    payment_blocking_reason TEXT,
    payment_methods_list TEXT,
    payment_terms TEXT,
    reconciliation_account TEXT NOT NULL DEFAULT '',
    deletion_indicator BOOLEAN NOT NULL DEFAULT FALSE,
    customer_account_group TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (customer, company_code),
    CONSTRAINT fk_cca_customer FOREIGN KEY (customer) REFERENCES business_partners(business_partner)
);

CREATE TABLE IF NOT EXISTS customer_sales_area_assignments (
    customer VARCHAR(10) NOT NULL,
    sales_organization VARCHAR(4) NOT NULL,
    distribution_channel VARCHAR(2) NOT NULL,
    division VARCHAR(2) NOT NULL,
    billing_is_blocked_for_customer TEXT,
    complete_delivery_is_defined BOOLEAN NOT NULL DEFAULT FALSE,
    credit_control_area TEXT,
    currency VARCHAR(3) NOT NULL DEFAULT '',
    customer_payment_terms VARCHAR(4) NOT NULL DEFAULT '',
    delivery_priority TEXT NOT NULL DEFAULT '',
    incoterms_classification VARCHAR(3) NOT NULL DEFAULT '',
    incoterms_location1 TEXT NOT NULL DEFAULT '',
    sales_group TEXT,
    sales_office TEXT,
    shipping_condition TEXT NOT NULL DEFAULT '',
    sls_unlmtd_ovrdeliv_is_allwd BOOLEAN NOT NULL DEFAULT FALSE,
    supplying_plant TEXT,
    sales_district TEXT,
    exchange_rate_type TEXT,
    PRIMARY KEY (customer, sales_organization, distribution_channel, division),
    CONSTRAINT fk_csaa_customer FOREIGN KEY (customer) REFERENCES business_partners(business_partner)
);

CREATE TABLE IF NOT EXISTS sales_order_headers (
    sales_order VARCHAR(10) PRIMARY KEY,
    sales_order_type VARCHAR(4) NOT NULL DEFAULT '',
    sales_organization VARCHAR(4) NOT NULL DEFAULT '',
    distribution_channel VARCHAR(2) NOT NULL DEFAULT '',
    organization_division VARCHAR(2) NOT NULL DEFAULT '',
    sales_group TEXT,
    sales_office TEXT,
    sold_to_party VARCHAR(10) NOT NULL DEFAULT '',
    creation_date TIMESTAMPTZ NOT NULL,
    created_by_user VARCHAR(12) NOT NULL DEFAULT '',
    last_change_date_time TIMESTAMPTZ NOT NULL,
    total_net_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    overall_delivery_status VARCHAR(1) NOT NULL DEFAULT '',
    overall_ord_reltd_billg_status TEXT,
    overall_sd_doc_reference_status TEXT,
    total_credit_check_status TEXT,
    transaction_currency VARCHAR(3) NOT NULL DEFAULT '',
    pricing_date TIMESTAMPTZ NOT NULL,
    requested_delivery_date TIMESTAMPTZ NOT NULL,
    header_billing_block_reason TEXT,
    delivery_block_reason TEXT,
    incoterms_classification VARCHAR(3) NOT NULL DEFAULT '',
    incoterms_location1 TEXT,
    customer_payment_terms VARCHAR(4) NOT NULL DEFAULT ''
);

CREATE INDEX idx_soh_sold_to_party ON sales_order_headers(sold_to_party);
CREATE INDEX idx_soh_creation_date ON sales_order_headers(creation_date);

CREATE TABLE IF NOT EXISTS sales_order_items (
    sales_order VARCHAR(10) NOT NULL,
    sales_order_item VARCHAR(6) NOT NULL,
    sales_order_item_category VARCHAR(4) NOT NULL DEFAULT '',
    material VARCHAR(40) NOT NULL DEFAULT '',
    requested_quantity TEXT NOT NULL DEFAULT '',
    requested_quantity_unit VARCHAR(3) NOT NULL DEFAULT '',
    transaction_currency VARCHAR(3) NOT NULL DEFAULT '',
    net_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    material_group VARCHAR(9) NOT NULL DEFAULT '',
    production_plant VARCHAR(4) NOT NULL DEFAULT '',
    storage_location VARCHAR(4) NOT NULL DEFAULT '',
    sales_document_rjcn_reason TEXT,
    item_billing_block_reason TEXT,
    PRIMARY KEY (sales_order, sales_order_item),
    CONSTRAINT fk_soi_header FOREIGN KEY (sales_order) REFERENCES sales_order_headers(sales_order)
);

CREATE INDEX idx_soi_material ON sales_order_items(material);
CREATE INDEX idx_soi_plant ON sales_order_items(production_plant);

CREATE TABLE IF NOT EXISTS sales_order_schedule_lines (
    sales_order VARCHAR(10) NOT NULL,
    sales_order_item VARCHAR(6) NOT NULL,
    schedule_line VARCHAR(4) NOT NULL,
    confirmed_delivery_date TIMESTAMPTZ NOT NULL,
    order_quantity_unit VARCHAR(3) NOT NULL DEFAULT '',
    confd_order_qty_by_matl_avail_check TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (sales_order, sales_order_item, schedule_line),
    CONSTRAINT fk_sosl_item FOREIGN KEY (sales_order, sales_order_item) REFERENCES sales_order_items(sales_order, sales_order_item)
);

CREATE TABLE IF NOT EXISTS outbound_delivery_headers (
    delivery_document VARCHAR(10) PRIMARY KEY,
    actual_goods_movement_date TIMESTAMPTZ,
    actual_goods_movement_time JSONB,
    creation_date TEXT NOT NULL DEFAULT '',
    creation_time JSONB,
    delivery_block_reason TEXT NOT NULL DEFAULT '',
    hdr_general_incompletion_status TEXT NOT NULL DEFAULT '',
    header_billing_block_reason TEXT NOT NULL DEFAULT '',
    last_change_date TIMESTAMPTZ,
    overall_goods_movement_status TEXT NOT NULL DEFAULT '',
    overall_picking_status TEXT NOT NULL DEFAULT '',
    overall_proof_of_delivery_status TEXT NOT NULL DEFAULT '',
    shipping_point TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS outbound_delivery_items (
    delivery_document VARCHAR(10) NOT NULL,
    delivery_document_item VARCHAR(6) NOT NULL,
    actual_delivery_quantity DECIMAL(15,3) NOT NULL DEFAULT 0,
    batch TEXT NOT NULL DEFAULT '',
    delivery_quantity_unit VARCHAR(3) NOT NULL DEFAULT '',
    item_billing_block_reason TEXT NOT NULL DEFAULT '',
    last_change_date TIMESTAMPTZ,
    plant VARCHAR(4) NOT NULL DEFAULT '',
    reference_sd_document TEXT NOT NULL DEFAULT '',
    reference_sd_document_item TEXT NOT NULL DEFAULT '',
    storage_location VARCHAR(4) NOT NULL DEFAULT '',
    PRIMARY KEY (delivery_document, delivery_document_item),
    CONSTRAINT fk_odi_header FOREIGN KEY (delivery_document) REFERENCES outbound_delivery_headers(delivery_document)
);

CREATE INDEX idx_odi_reference_sd ON outbound_delivery_items(reference_sd_document);
CREATE INDEX idx_odi_plant ON outbound_delivery_items(plant);

CREATE TABLE IF NOT EXISTS billing_document_headers (
    billing_document VARCHAR(10) PRIMARY KEY,
    billing_document_type VARCHAR(4) NOT NULL DEFAULT '',
    company_code VARCHAR(4) NOT NULL DEFAULT '',
    fiscal_year VARCHAR(4) NOT NULL DEFAULT '',
    creation_date TIMESTAMPTZ NOT NULL,
    creation_time JSONB,
    last_change_date_time TIMESTAMPTZ,
    billing_document_date TIMESTAMPTZ NOT NULL,
    billing_document_is_cancelled BOOLEAN,
    cancelled_billing_document TEXT,
    total_net_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    transaction_currency VARCHAR(3) NOT NULL DEFAULT '',
    accounting_document TEXT,
    sold_to_party VARCHAR(10) NOT NULL DEFAULT ''
);

CREATE INDEX idx_bdh_sold_to_party ON billing_document_headers(sold_to_party);

CREATE TABLE IF NOT EXISTS billing_document_items (
    billing_document VARCHAR(10) NOT NULL,
    billing_document_item VARCHAR(6) NOT NULL,
    material VARCHAR(40) NOT NULL DEFAULT '',
    billing_quantity TEXT NOT NULL DEFAULT '',
    billing_quantity_unit VARCHAR(3) NOT NULL DEFAULT '',
    transaction_currency VARCHAR(3) NOT NULL DEFAULT '',
    net_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    reference_sd_document TEXT,
    reference_sd_document_item TEXT,
    PRIMARY KEY (billing_document, billing_document_item),
    CONSTRAINT fk_bdi_header FOREIGN KEY (billing_document) REFERENCES billing_document_headers(billing_document)
);

CREATE INDEX idx_bdi_reference_sd ON billing_document_items(reference_sd_document);

CREATE TABLE IF NOT EXISTS billing_document_cancellations (
    billing_document VARCHAR(10) PRIMARY KEY,
    billing_document_type VARCHAR(4) NOT NULL DEFAULT '',
    company_code VARCHAR(4) NOT NULL DEFAULT '',
    fiscal_year VARCHAR(4) NOT NULL DEFAULT '',
    creation_date TIMESTAMPTZ NOT NULL,
    creation_time JSONB,
    last_change_date_time TIMESTAMPTZ,
    billing_document_date TIMESTAMPTZ NOT NULL,
    billing_document_is_cancelled BOOLEAN,
    cancelled_billing_document TEXT,
    total_net_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    transaction_currency VARCHAR(3) NOT NULL DEFAULT '',
    accounting_document TEXT,
    sold_to_party VARCHAR(10) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS journal_entry_items_accounts_receivable (
    company_code VARCHAR(4) NOT NULL,
    fiscal_year INTEGER NOT NULL,
    accounting_document VARCHAR(10) NOT NULL,
    gl_account VARCHAR(10) NOT NULL,
    reference_document VARCHAR(10) NOT NULL,
    cost_center TEXT,
    profit_center VARCHAR(10) NOT NULL,
    transaction_currency VARCHAR(3) NOT NULL,
    amount_in_transaction_currency DECIMAL(20,4) NOT NULL DEFAULT 0,
    company_code_currency VARCHAR(3) NOT NULL,
    amount_in_company_code_currency DECIMAL(20,4) NOT NULL DEFAULT 0,
    posting_date TIMESTAMPTZ,
    document_date TIMESTAMPTZ,
    accounting_document_type VARCHAR(2) NOT NULL,
    accounting_document_item VARCHAR(3) NOT NULL,
    assignment_reference TEXT,
    last_change_date_time TIMESTAMPTZ NOT NULL,
    customer VARCHAR(10) NOT NULL,
    financial_account_type VARCHAR(1) NOT NULL,
    clearing_date TIMESTAMPTZ,
    clearing_accounting_document VARCHAR(10) NOT NULL,
    clearing_doc_fiscal_year INTEGER,
    PRIMARY KEY (fiscal_year, accounting_document, accounting_document_item)
);

CREATE INDEX idx_jeiar_customer ON journal_entry_items_accounts_receivable(customer);
CREATE INDEX idx_jeiar_reference_document ON journal_entry_items_accounts_receivable(reference_document);
CREATE INDEX idx_jeiar_clearing_doc ON journal_entry_items_accounts_receivable(clearing_accounting_document);

CREATE TABLE IF NOT EXISTS payments_accounts_receivable (
    company_code VARCHAR(4) NOT NULL,
    fiscal_year INTEGER NOT NULL,
    accounting_document VARCHAR(10) NOT NULL,
    accounting_document_item VARCHAR(3) NOT NULL,
    clearing_date TIMESTAMPTZ,
    clearing_accounting_document VARCHAR(10) NOT NULL,
    clearing_doc_fiscal_year INTEGER,
    amount_in_transaction_currency DECIMAL(20,4) NOT NULL DEFAULT 0,
    transaction_currency VARCHAR(3) NOT NULL,
    amount_in_company_code_currency DECIMAL(20,4) NOT NULL DEFAULT 0,
    company_code_currency VARCHAR(3) NOT NULL,
    customer VARCHAR(10) NOT NULL,
    invoice_reference TEXT,
    invoice_reference_fiscal_year INTEGER,
    sales_document TEXT,
    sales_document_item TEXT,
    posting_date TIMESTAMPTZ,
    document_date TIMESTAMPTZ,
    assignment_reference TEXT,
    gl_account VARCHAR(10) NOT NULL,
    financial_account_type VARCHAR(1) NOT NULL,
    profit_center VARCHAR(10) NOT NULL,
    cost_center TEXT,
    PRIMARY KEY (fiscal_year, accounting_document, accounting_document_item)
);

CREATE INDEX idx_par_customer ON payments_accounts_receivable(customer);
CREATE INDEX idx_par_clearing_doc ON payments_accounts_receivable(clearing_accounting_document);
