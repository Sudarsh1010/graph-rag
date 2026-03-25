package loader

import (
	"context"

	"github.com/sudarsh1010/graph-rag/internal/model"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type BillingDocumentHeadersLoader struct{ baseLoader }

func NewBillingDocumentHeadersLoader() *BillingDocumentHeadersLoader {
	return &BillingDocumentHeadersLoader{baseLoader: newBaseLoader("billing_document_headers", "billing_document_headers", "billing_document_headers")}
}

func (l *BillingDocumentHeadersLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *BillingDocumentHeadersLoader) convert(raw map[string]interface{}) (interface{}, error) {
	creationDate, _ := getTime(raw, "creationDate")
	creationTime := mustTimeOfDay(raw, "creationTime")
	lastChangeDateTime, _ := getTimePtr(raw, "lastChangeDateTime")
	billingDate, _ := getTime(raw, "billingDocumentDate")
	totalNetAmount, _ := getDecimal(raw, "totalNetAmount")
	isCancelled := getBool(raw, "billingDocumentIsCancelled")

	m := &model.BillingDocumentHeader{
		BillingDocument:            getString(raw, "billingDocument"),
		BillingDocumentType:        getString(raw, "billingDocumentType"),
		CompanyCode:                getString(raw, "companyCode"),
		FiscalYear:                 getString(raw, "fiscalYear"),
		CreationDate:               creationDate,
		CreationTime:               creationTime,
		LastChangeDateTime:         lastChangeDateTime,
		BillingDocumentDate:        billingDate,
		BillingDocumentIsCancelled: &isCancelled,
		CancelledBillingDocument:   getStringPtr(raw, "cancelledBillingDocument"),
		TotalNetAmount:             totalNetAmount,
		TransactionCurrency:        getString(raw, "transactionCurrency"),
		AccountingDocument:         getStringPtr(raw, "accountingDocument"),
		SoldToParty:                getString(raw, "soldToParty"),
	}
	if err := ValidateRequired(m.BillingDocument, "billingDocument"); err != nil {
		return nil, err
	}
	return m, nil
}

type BillingDocumentItemsLoader struct{ baseLoader }

func NewBillingDocumentItemsLoader() *BillingDocumentItemsLoader {
	return &BillingDocumentItemsLoader{baseLoader: newBaseLoader("billing_document_items", "billing_document_items", "billing_document_items")}
}

func (l *BillingDocumentItemsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *BillingDocumentItemsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	netAmount, _ := getDecimal(raw, "netAmount")

	m := &model.BillingDocumentItem{
		BillingDocument:         getString(raw, "billingDocument"),
		BillingDocumentItem:     getString(raw, "billingDocumentItem"),
		Material:                getString(raw, "material"),
		BillingQuantity:         getString(raw, "billingQuantity"),
		BillingQuantityUnit:     getString(raw, "billingQuantityUnit"),
		TransactionCurrency:     getString(raw, "transactionCurrency"),
		NetAmount:               netAmount,
		ReferenceSdDocument:     getStringPtr(raw, "referenceSdDocument"),
		ReferenceSdDocumentItem: getStringPtr(raw, "referenceSdDocumentItem"),
	}
	return m, nil
}

type BillingDocumentCancellationsLoader struct{ baseLoader }

func NewBillingDocumentCancellationsLoader() *BillingDocumentCancellationsLoader {
	return &BillingDocumentCancellationsLoader{baseLoader: newBaseLoader("billing_document_cancellations", "billing_document_cancellations", "billing_document_cancellations")}
}

func (l *BillingDocumentCancellationsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *BillingDocumentCancellationsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	creationDate, _ := getTime(raw, "creationDate")
	creationTime := mustTimeOfDay(raw, "creationTime")
	lastChangeDateTime, _ := getTimePtr(raw, "lastChangeDateTime")
	billingDate, _ := getTime(raw, "billingDocumentDate")
	totalNetAmount, _ := getDecimal(raw, "totalNetAmount")
	isCancelled := getBool(raw, "billingDocumentIsCancelled")

	m := &model.BillingDocumentCancellation{
		BillingDocument:            getString(raw, "billingDocument"),
		BillingDocumentType:        getString(raw, "billingDocumentType"),
		CompanyCode:                getString(raw, "companyCode"),
		FiscalYear:                 getString(raw, "fiscalYear"),
		CreationDate:               creationDate,
		CreationTime:               creationTime,
		LastChangeDateTime:         lastChangeDateTime,
		BillingDocumentDate:        billingDate,
		BillingDocumentIsCancelled: &isCancelled,
		CancelledBillingDocument:   getStringPtr(raw, "cancelledBillingDocument"),
		TotalNetAmount:             totalNetAmount,
		TransactionCurrency:        getString(raw, "transactionCurrency"),
		AccountingDocument:         getStringPtr(raw, "accountingDocument"),
		SoldToParty:                getString(raw, "soldToParty"),
	}
	if err := ValidateRequired(m.BillingDocument, "billingDocument"); err != nil {
		return nil, err
	}
	return m, nil
}
