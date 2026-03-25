package loader

import (
	"context"

	"github.com/sudarsh1010/graph-rag/internal/model"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type JournalEntryItemsARLoader struct{ baseLoader }

func NewJournalEntryItemsARLoader() *JournalEntryItemsARLoader {
	return &JournalEntryItemsARLoader{baseLoader: newBaseLoader("journal_entry_items_ar", "journal_entry_items_accounts_receivable", "journal_entry_items_accounts_receivable")}
}

func (l *JournalEntryItemsARLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *JournalEntryItemsARLoader) convert(raw map[string]interface{}) (interface{}, error) {
	fiscalYear, err := getIntRequired(raw, "fiscalYear")
	if err != nil {
		return nil, err
	}
	amtTxn, _ := getDecimal(raw, "amountInTransactionCurrency")
	amtCoCode, _ := getDecimal(raw, "amountInCompanyCodeCurrency")
	postingDate, _ := getTimePtr(raw, "postingDate")
	documentDate, _ := getTimePtr(raw, "documentDate")
	lastChangeDateTime, _ := getTime(raw, "lastChangeDateTime")
	clearingDate, _ := getTimePtr(raw, "clearingDate")
	clearingDocFiscalYear, _ := getIntPtr(raw, "clearingDocFiscalYear")

	m := &model.JournalEntryItemAR{
		CompanyCode:                 getString(raw, "companyCode"),
		FiscalYear:                  fiscalYear,
		AccountingDocument:          getString(raw, "accountingDocument"),
		GLAccount:                   getString(raw, "glAccount"),
		ReferenceDocument:           getString(raw, "referenceDocument"),
		CostCenter:                  getStringPtr(raw, "costCenter"),
		ProfitCenter:                getString(raw, "profitCenter"),
		TransactionCurrency:         getString(raw, "transactionCurrency"),
		AmountInTransactionCurrency: amtTxn,
		CompanyCodeCurrency:         getString(raw, "companyCodeCurrency"),
		AmountInCompanyCodeCurrency: amtCoCode,
		PostingDate:                 postingDate,
		DocumentDate:                documentDate,
		AccountingDocumentType:      getString(raw, "accountingDocumentType"),
		AccountingDocumentItem:      getString(raw, "accountingDocumentItem"),
		AssignmentReference:         getStringPtr(raw, "assignmentReference"),
		LastChangeDateTime:          lastChangeDateTime,
		Customer:                    getString(raw, "customer"),
		FinancialAccountType:        getString(raw, "financialAccountType"),
		ClearingDate:                clearingDate,
		ClearingAccountingDocument:  getString(raw, "clearingAccountingDocument"),
		ClearingDocFiscalYear:       clearingDocFiscalYear,
	}
	return m, nil
}

type PaymentsARLoader struct{ baseLoader }

func NewPaymentsARLoader() *PaymentsARLoader {
	return &PaymentsARLoader{baseLoader: newBaseLoader("payments_ar", "payments_accounts_receivable", "payments_accounts_receivable")}
}

func (l *PaymentsARLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *PaymentsARLoader) convert(raw map[string]interface{}) (interface{}, error) {
	fiscalYear, err := getIntRequired(raw, "fiscalYear")
	if err != nil {
		return nil, err
	}
	amtTxn, _ := getDecimal(raw, "amountInTransactionCurrency")
	amtCoCode, _ := getDecimal(raw, "amountInCompanyCodeCurrency")
	clearingDate, _ := getTimePtr(raw, "clearingDate")
	clearingDocFiscalYear, _ := getIntPtr(raw, "clearingDocFiscalYear")
	postingDate, _ := getTimePtr(raw, "postingDate")
	documentDate, _ := getTimePtr(raw, "documentDate")

	m := &model.PaymentAR{
		CompanyCode:                 getString(raw, "companyCode"),
		FiscalYear:                  fiscalYear,
		AccountingDocument:          getString(raw, "accountingDocument"),
		AccountingDocumentItem:      getString(raw, "accountingDocumentItem"),
		ClearingDate:                clearingDate,
		ClearingAccountingDocument:  getString(raw, "clearingAccountingDocument"),
		ClearingDocFiscalYear:       clearingDocFiscalYear,
		AmountInTransactionCurrency: amtTxn,
		TransactionCurrency:         getString(raw, "transactionCurrency"),
		AmountInCompanyCodeCurrency: amtCoCode,
		CompanyCodeCurrency:         getString(raw, "companyCodeCurrency"),
		Customer:                    getString(raw, "customer"),
		InvoiceReference:            getStringPtr(raw, "invoiceReference"),
		InvoiceReferenceFiscalYear:  mustIntPtr(raw, "invoiceReferenceFiscalYear"),
		SalesDocument:               getStringPtr(raw, "salesDocument"),
		SalesDocumentItem:           getStringPtr(raw, "salesDocumentItem"),
		PostingDate:                 postingDate,
		DocumentDate:                documentDate,
		AssignmentReference:         getStringPtr(raw, "assignmentReference"),
		GLAccount:                   getString(raw, "glAccount"),
		FinancialAccountType:        getString(raw, "financialAccountType"),
		ProfitCenter:                getString(raw, "profitCenter"),
		CostCenter:                  getStringPtr(raw, "costCenter"),
	}
	return m, nil
}
