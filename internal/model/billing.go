package model

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

// BillingDocumentHeader represents a billing document header in the system.
// Primary key: billingDocument
type BillingDocumentHeader struct {
	bun.BaseModel `bun:"table:billing_document_headers"`

	// Primary Key
	BillingDocument string `bun:",pk" json:"billingDocument"`

	// Order Information
	BillingDocumentType string `json:"billingDocumentType"`
	CompanyCode         string `json:"companyCode"`
	FiscalYear          string `json:"fiscalYear"`

	// Timestamps
	CreationDate       time.Time  `json:"creationDate"`
	CreationTime       TimeOfDay  `json:"creationTime"`
	LastChangeDateTime *time.Time `bun:",nullzero" json:"lastChangeDateTime,omitempty"`

	// Dates
	BillingDocumentDate time.Time `json:"billingDocumentDate"`

	// Status Flags
	BillingDocumentIsCancelled *bool   `bun:",nullzero" json:"billingDocumentIsCancelled,omitempty"`
	CancelledBillingDocument   *string `bun:",nullzero" json:"cancelledBillingDocument,omitempty"`

	// Amounts
	TotalNetAmount decimal.Decimal `json:"totalNetAmount"`

	// Other Fields
	TransactionCurrency string  `json:"transactionCurrency"`
	AccountingDocument  *string `bun:",nullzero" json:"accountingDocument,omitempty"`
	SoldToParty         string  `json:"soldToParty"`
}

// BillingDocumentItem represents an item within a billing document.
// Composite primary key: billingDocument + billingDocumentItem
type BillingDocumentItem struct {
	bun.BaseModel `bun:"table:billing_document_items"`

	// Composite Primary Key
	BillingDocument     string `bun:",pk" json:"billingDocument"`
	BillingDocumentItem string `bun:",pk" json:"billingDocumentItem"`

	// Item Information
	Material            string          `json:"material"`
	BillingQuantity     string          `json:"billingQuantity"`
	BillingQuantityUnit string          `json:"billingQuantityUnit"`
	TransactionCurrency string          `json:"transactionCurrency"`
	NetAmount           decimal.Decimal `json:"netAmount"`

	// Reference Information
	ReferenceSdDocument     *string `bun:",nullzero" json:"referenceSdDocument,omitempty"`
	ReferenceSdDocumentItem *string `bun:",nullzero" json:"referenceSdDocumentItem,omitempty"`
}

// BillingDocumentCancellation represents a billing document cancellation.
// Primary key: billingDocument
type BillingDocumentCancellation struct {
	bun.BaseModel `bun:"table:billing_document_cancellations"`

	// Primary Key
	BillingDocument string `bun:",pk" json:"billingDocument"`

	// Order Information
	BillingDocumentType string `json:"billingDocumentType"`
	CompanyCode         string `json:"companyCode"`
	FiscalYear          string `json:"fiscalYear"`

	// Timestamps
	CreationDate       time.Time  `json:"creationDate"`
	CreationTime       TimeOfDay  `json:"creationTime"`
	LastChangeDateTime *time.Time `bun:",nullzero" json:"lastChangeDateTime,omitempty"`

	// Dates
	BillingDocumentDate time.Time `json:"billingDocumentDate"`

	// Status Flags
	BillingDocumentIsCancelled *bool   `bun:",nullzero" json:"billingDocumentIsCancelled,omitempty"`
	CancelledBillingDocument   *string `bun:",nullzero" json:"cancelledBillingDocument,omitempty"`

	// Amounts
	TotalNetAmount decimal.Decimal `json:"totalNetAmount"`

	// Other Fields
	TransactionCurrency string  `json:"transactionCurrency"`
	AccountingDocument  *string `bun:",nullzero" json:"accountingDocument,omitempty"`
	SoldToParty         string  `json:"soldToParty"`
}
