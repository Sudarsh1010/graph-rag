package model

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// SalesOrderHeader represents a sales order header in the system.
// Primary key: salesOrder
type SalesOrderHeader struct {
	bun.BaseModel `bun:"table:sales_order_headers"`

	// Primary Key
	SalesOrder string `bun:",pk" json:"salesOrder"`

	// Order Information
	SalesOrderType       string `json:"salesOrderType"`
	SalesOrganization    string `json:"salesOrganization"`
	DistributionChannel  string `json:"distributionChannel"`
	OrganizationDivision string `json:"organizationDivision"`

	// Optional Fields
	SalesGroup  *string `bun:",nullzero" json:"salesGroup,omitempty"`
	SalesOffice *string `bun:",nullzero" json:"salesOffice,omitempty"`

	// Customer Information
	SoldToParty string `json:"soldToParty"`

	// Timestamps
	CreationDate       time.Time `json:"creationDate"`
	CreatedByUser      string    `json:"createdByUser"`
	LastChangeDateTime time.Time `json:"lastChangeDateTime"`

	// Amounts
	TotalNetAmount decimal.Decimal `json:"totalNetAmount"`

	// Status Flags
	OverallDeliveryStatus       string  `json:"overallDeliveryStatus"`
	OverallOrdReltdBillgStatus  *string `bun:",nullzero" json:"overallOrdReltdBillgStatus,omitempty"`
	OverallSdDocReferenceStatus *string `bun:",nullzero" json:"overallSdDocReferenceStatus,omitempty"`
	TotalCreditCheckStatus      *string `bun:",nullzero" json:"totalCreditCheckStatus,omitempty"`

	// Order Details
	TransactionCurrency   string    `json:"transactionCurrency"`
	PricingDate           time.Time `json:"pricingDate"`
	RequestedDeliveryDate time.Time `json:"requestedDeliveryDate"`

	// Blocking Reasons
	HeaderBillingBlockReason *string `bun:",nullzero" json:"headerBillingBlockReason,omitempty"`
	DeliveryBlockReason      *string `bun:",nullzero" json:"deliveryBlockReason,omitempty"`

	// Incoterms
	IncotermsClassification string  `json:"incotermsClassification"`
	IncotermsLocation1      *string `bun:",nullzero" json:"incotermsLocation1,omitempty"`

	// Payment
	CustomerPaymentTerms string `json:"customerPaymentTerms"`
}

// SalesOrderItem represents an item within a sales order.
// Composite primary key: salesOrder + salesOrderItem
type SalesOrderItem struct {
	bun.BaseModel `bun:"table:sales_order_items"`

	// Composite Primary Key
	SalesOrder     string `bun:",pk" json:"salesOrder"`
	SalesOrderItem string `bun:",pk" json:"salesOrderItem"`

	// Item Information
	SalesOrderItemCategory string `json:"salesOrderItemCategory"`
	Material               string `json:"material"`

	// Quantity and Currency
	RequestedQuantity     string          `json:"requestedQuantity"`
	RequestedQuantityUnit string          `json:"requestedQuantityUnit"`
	TransactionCurrency   string          `json:"transactionCurrency"`
	NetAmount             decimal.Decimal `json:"netAmount"`

	// Material Information
	MaterialGroup   string `json:"materialGroup"`
	ProductionPlant string `json:"productionPlant"`
	StorageLocation string `json:"storageLocation"`

	// Optional Fields
	SalesDocumentRjcnReason *string `bun:",nullzero" json:"salesDocumentRjcnReason,omitempty"`
	ItemBillingBlockReason  *string `bun:",nullzero" json:"itemBillingBlockReason,omitempty"`
}

// SalesOrderScheduleLine represents a schedule line for a sales order item.
// Composite primary key: salesOrder + salesOrderItem + scheduleLine
type SalesOrderScheduleLine struct {
	bun.BaseModel `bun:"table:sales_order_schedule_lines"`

	// Composite Primary Key
	SalesOrder     string `bun:",pk" json:"salesOrder"`
	SalesOrderItem string `bun:",pk" json:"salesOrderItem"`
	ScheduleLine   string `bun:",pk" json:"scheduleLine"`

	// Schedule Information
	ConfirmedDeliveryDate         time.Time `json:"confirmedDeliveryDate"`
	OrderQuantityUnit             string    `json:"orderQuantityUnit"`
	ConfdOrderQtyByMatlAvailCheck string    `json:"confdOrderQtyByMatlAvailCheck"`
}
