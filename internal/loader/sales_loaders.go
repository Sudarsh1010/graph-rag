package loader

import (
	"context"

	"github.com/sudarsh1010/graph-rag/internal/model"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type SalesOrderHeadersLoader struct{ baseLoader }

func NewSalesOrderHeadersLoader() *SalesOrderHeadersLoader {
	return &SalesOrderHeadersLoader{baseLoader: newBaseLoader("sales_order_headers", "sales_order_headers", "sales_order_headers")}
}

func (l *SalesOrderHeadersLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *SalesOrderHeadersLoader) convert(raw map[string]interface{}) (interface{}, error) {
	creationDate, _ := getTime(raw, "creationDate")
	createdByUser := getString(raw, "createdByUser")
	lastChangeDateTime, _ := getTime(raw, "lastChangeDateTime")
	totalNetAmount, _ := getDecimal(raw, "totalNetAmount")
	pricingDate, _ := getTime(raw, "pricingDate")
	requestedDeliveryDate, _ := getTime(raw, "requestedDeliveryDate")

	m := &model.SalesOrderHeader{
		SalesOrder:                  getString(raw, "salesOrder"),
		SalesOrderType:              getString(raw, "salesOrderType"),
		SalesOrganization:           getString(raw, "salesOrganization"),
		DistributionChannel:         getString(raw, "distributionChannel"),
		OrganizationDivision:        getString(raw, "organizationDivision"),
		SalesGroup:                  getStringPtr(raw, "salesGroup"),
		SalesOffice:                 getStringPtr(raw, "salesOffice"),
		SoldToParty:                 getString(raw, "soldToParty"),
		CreationDate:                creationDate,
		CreatedByUser:               createdByUser,
		LastChangeDateTime:          lastChangeDateTime,
		TotalNetAmount:              totalNetAmount,
		OverallDeliveryStatus:       getString(raw, "overallDeliveryStatus"),
		OverallOrdReltdBillgStatus:  getStringPtr(raw, "overallOrdReltdBillgStatus"),
		OverallSdDocReferenceStatus: getStringPtr(raw, "overallSdDocReferenceStatus"),
		TotalCreditCheckStatus:      getStringPtr(raw, "totalCreditCheckStatus"),
		TransactionCurrency:         getString(raw, "transactionCurrency"),
		PricingDate:                 pricingDate,
		RequestedDeliveryDate:       requestedDeliveryDate,
		HeaderBillingBlockReason:    getStringPtr(raw, "headerBillingBlockReason"),
		DeliveryBlockReason:         getStringPtr(raw, "deliveryBlockReason"),
		IncotermsClassification:     getString(raw, "incotermsClassification"),
		IncotermsLocation1:          getStringPtr(raw, "incotermsLocation1"),
		CustomerPaymentTerms:        getString(raw, "customerPaymentTerms"),
	}
	if err := ValidateRequired(m.SalesOrder, "salesOrder"); err != nil {
		return nil, err
	}
	return m, nil
}

type SalesOrderItemsLoader struct{ baseLoader }

func NewSalesOrderItemsLoader() *SalesOrderItemsLoader {
	return &SalesOrderItemsLoader{baseLoader: newBaseLoader("sales_order_items", "sales_order_items", "sales_order_items")}
}

func (l *SalesOrderItemsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *SalesOrderItemsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	netAmount, _ := getDecimal(raw, "netAmount")

	m := &model.SalesOrderItem{
		SalesOrder:              getString(raw, "salesOrder"),
		SalesOrderItem:          getString(raw, "salesOrderItem"),
		SalesOrderItemCategory:  getString(raw, "salesOrderItemCategory"),
		Material:                getString(raw, "material"),
		RequestedQuantity:       getString(raw, "requestedQuantity"),
		RequestedQuantityUnit:   getString(raw, "requestedQuantityUnit"),
		TransactionCurrency:     getString(raw, "transactionCurrency"),
		NetAmount:               netAmount,
		MaterialGroup:           getString(raw, "materialGroup"),
		ProductionPlant:         getString(raw, "productionPlant"),
		StorageLocation:         getString(raw, "storageLocation"),
		SalesDocumentRjcnReason: getStringPtr(raw, "salesDocumentRjcnReason"),
		ItemBillingBlockReason:  getStringPtr(raw, "itemBillingBlockReason"),
	}
	return m, nil
}

type SalesOrderScheduleLinesLoader struct{ baseLoader }

func NewSalesOrderScheduleLinesLoader() *SalesOrderScheduleLinesLoader {
	return &SalesOrderScheduleLinesLoader{baseLoader: newBaseLoader("sales_order_schedule_lines", "sales_order_schedule_lines", "sales_order_schedule_lines")}
}

func (l *SalesOrderScheduleLinesLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *SalesOrderScheduleLinesLoader) convert(raw map[string]interface{}) (interface{}, error) {
	confirmedDeliveryDate, _ := getTime(raw, "confirmedDeliveryDate")

	m := &model.SalesOrderScheduleLine{
		SalesOrder:                    getString(raw, "salesOrder"),
		SalesOrderItem:                getString(raw, "salesOrderItem"),
		ScheduleLine:                  getString(raw, "scheduleLine"),
		ConfirmedDeliveryDate:         confirmedDeliveryDate,
		OrderQuantityUnit:             getString(raw, "orderQuantityUnit"),
		ConfdOrderQtyByMatlAvailCheck: getString(raw, "confdOrderQtyByMatlAvailCheck"),
	}
	return m, nil
}
