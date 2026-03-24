package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type JournalEntryItemAR struct {
	CompanyCode                 string          `bun:"companyCode,notnull" json:"companyCode,omitempty"`
	FiscalYear                  int             `bun:"fiscalYear,notnull,pk" json:"fiscalYear,omitempty"`
	AccountingDocument          string          `bun:"accountingDocument,notnull,pk" json:"accountingDocument,omitempty"`
	GLAccount                   string          `bun:"glAccount,notnull" json:"glAccount,omitempty"`
	ReferenceDocument           string          `bun:"referenceDocument,notnull" json:"referenceDocument,omitempty"`
	CostCenter                  *string         `bun:"costCenter,nullzero" json:"costCenter,omitempty"`
	ProfitCenter                string          `bun:"profitCenter,notnull" json:"profitCenter,omitempty"`
	TransactionCurrency         string          `bun:"transactionCurrency,notnull" json:"transactionCurrency,omitempty"`
	AmountInTransactionCurrency decimal.Decimal `bun:"amountInTransactionCurrency,notnull,type:decimal(20,4)" json:"amountInTransactionCurrency,omitempty"`
	CompanyCodeCurrency         string          `bun:"companyCodeCurrency,notnull" json:"companyCodeCurrency,omitempty"`
	AmountInCompanyCodeCurrency decimal.Decimal `bun:"amountInCompanyCodeCurrency,notnull,type:decimal(20,4)" json:"amountInCompanyCodeCurrency,omitempty"`
	PostingDate                 *time.Time      `bun:"postingDate,nullzero" json:"postingDate,omitempty"`
	DocumentDate                *time.Time      `bun:"documentDate,nullzero" json:"documentDate,omitempty"`
	AccountingDocumentType      string          `bun:"accountingDocumentType,notnull" json:"accountingDocumentType,omitempty"`
	AccountingDocumentItem      string          `bun:"accountingDocumentItem,notnull,pk" json:"accountingDocumentItem,omitempty"`
	AssignmentReference         *string         `bun:"assignmentReference,nullzero" json:"assignmentReference,omitempty"`
	LastChangeDateTime          time.Time       `bun:"lastChangeDateTime,notnull" json:"lastChangeDateTime,omitempty"`
	Customer                    string          `bun:"customer,notnull" json:"customer,omitempty"`
	FinancialAccountType        string          `bun:"financialAccountType,notnull" json:"financialAccountType,omitempty"`
	ClearingDate                *time.Time      `bun:"clearingDate,nullzero" json:"clearingDate,omitempty"`
	ClearingAccountingDocument  string          `bun:"clearingAccountingDocument,notnull" json:"clearingAccountingDocument,omitempty"`
	ClearingDocFiscalYear       *int            `bun:"clearingDocFiscalYear,nullzero" json:"clearingDocFiscalYear,omitempty"`
}

func (JournalEntryItemAR) TableName() string {
	return "journal_entry_items_accounts_receivable"
}

type PaymentAR struct {
	CompanyCode                 string          `bun:"companyCode,notnull" json:"companyCode,omitempty"`
	FiscalYear                  int             `bun:"fiscalYear,notnull,pk" json:"fiscalYear,omitempty"`
	AccountingDocument          string          `bun:"accountingDocument,notnull,pk" json:"accountingDocument,omitempty"`
	AccountingDocumentItem      string          `bun:"accountingDocumentItem,notnull,pk" json:"accountingDocumentItem,omitempty"`
	ClearingDate                *time.Time      `bun:"clearingDate,nullzero" json:"clearingDate,omitempty"`
	ClearingAccountingDocument  string          `bun:"clearingAccountingDocument,notnull" json:"clearingAccountingDocument,omitempty"`
	ClearingDocFiscalYear       *int            `bun:"clearingDocFiscalYear,nullzero" json:"clearingDocFiscalYear,omitempty"`
	AmountInTransactionCurrency decimal.Decimal `bun:"amountInTransactionCurrency,notnull,type:decimal(20,4)" json:"amountInTransactionCurrency,omitempty"`
	TransactionCurrency         string          `bun:"transactionCurrency,notnull" json:"transactionCurrency,omitempty"`
	AmountInCompanyCodeCurrency decimal.Decimal `bun:"amountInCompanyCodeCurrency,notnull,type:decimal(20,4)" json:"amountInCompanyCodeCurrency,omitempty"`
	CompanyCodeCurrency         string          `bun:"companyCodeCurrency,notnull" json:"companyCodeCurrency,omitempty"`
	Customer                    string          `bun:"customer,notnull" json:"customer,omitempty"`
	InvoiceReference            *string         `bun:"invoiceReference,nullzero" json:"invoiceReference,omitempty"`
	InvoiceReferenceFiscalYear  *int            `bun:"invoiceReferenceFiscalYear,nullzero" json:"invoiceReferenceFiscalYear,omitempty"`
	SalesDocument               *string         `bun:"salesDocument,nullzero" json:"salesDocument,omitempty"`
	SalesDocumentItem           *string         `bun:"salesDocumentItem,nullzero" json:"salesDocumentItem,omitempty"`
	PostingDate                 *time.Time      `bun:"postingDate,nullzero" json:"postingDate,omitempty"`
	DocumentDate                *time.Time      `bun:"documentDate,nullzero" json:"documentDate,omitempty"`
	AssignmentReference         *string         `bun:"assignmentReference,nullzero" json:"assignmentReference,omitempty"`
	GLAccount                   string          `bun:"glAccount,notnull" json:"glAccount,omitempty"`
	FinancialAccountType        string          `bun:"financialAccountType,notnull" json:"financialAccountType,omitempty"`
	ProfitCenter                string          `bun:"profitCenter,notnull" json:"profitCenter,omitempty"`
	CostCenter                  *string         `bun:"costCenter,nullzero" json:"costCenter,omitempty"`
}

func (PaymentAR) TableName() string {
	return "payments_accounts_receivable"
}
