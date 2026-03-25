package loader

import (
	"context"

	"github.com/sudarsh1010/graph-rag/internal/model"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type OutboundDeliveryHeadersLoader struct{ baseLoader }

func NewOutboundDeliveryHeadersLoader() *OutboundDeliveryHeadersLoader {
	return &OutboundDeliveryHeadersLoader{baseLoader: newBaseLoader("outbound_delivery_headers", "outbound_delivery_headers", "outbound_delivery_headers")}
}

func (l *OutboundDeliveryHeadersLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *OutboundDeliveryHeadersLoader) convert(raw map[string]interface{}) (interface{}, error) {
	actualGoodsMovementDate, _ := getTimePtr(raw, "actualGoodsMovementDate")
	actualGoodsMovementTime := mustTimeOfDay(raw, "actualGoodsMovementTime")
	lastChangeDate, _ := getTimePtr(raw, "lastChangeDate")
	creationTime := mustTimeOfDay(raw, "creationTime")

	m := &model.OutboundDeliveryHeader{
		ActualGoodsMovementDate:      actualGoodsMovementDate,
		ActualGoodsMovementTime:      actualGoodsMovementTime,
		CreationDate:                 getString(raw, "creationDate"),
		CreationTime:                 creationTime,
		DeliveryBlockReason:          getString(raw, "deliveryBlockReason"),
		DeliveryDocument:             getString(raw, "deliveryDocument"),
		HdrGeneralIncompletionStatus: getString(raw, "hdrGeneralIncompletionStatus"),
		HeaderBillingBlockReason:     getString(raw, "headerBillingBlockReason"),
		LastChangeDate:               lastChangeDate,
		OverallGoodsMovementStatus:   getString(raw, "overallGoodsMovementStatus"),
		OverallPickingStatus:         getString(raw, "overallPickingStatus"),
		OverallProofOfDeliveryStatus: getString(raw, "overallProofOfDeliveryStatus"),
		ShippingPoint:                getString(raw, "shippingPoint"),
	}
	if err := ValidateRequired(m.DeliveryDocument, "deliveryDocument"); err != nil {
		return nil, err
	}
	return m, nil
}

type OutboundDeliveryItemsLoader struct{ baseLoader }

func NewOutboundDeliveryItemsLoader() *OutboundDeliveryItemsLoader {
	return &OutboundDeliveryItemsLoader{baseLoader: newBaseLoader("outbound_delivery_items", "outbound_delivery_items", "outbound_delivery_items")}
}

func (l *OutboundDeliveryItemsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *OutboundDeliveryItemsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	actualDeliveryQuantity, _ := getDecimal(raw, "actualDeliveryQuantity")
	lastChangeDate, _ := getTimePtr(raw, "lastChangeDate")

	m := &model.OutboundDeliveryItem{
		ActualDeliveryQuantity:  actualDeliveryQuantity,
		Batch:                   getString(raw, "batch"),
		DeliveryDocument:        getString(raw, "deliveryDocument"),
		DeliveryDocumentItem:    getString(raw, "deliveryDocumentItem"),
		DeliveryQuantityUnit:    getString(raw, "deliveryQuantityUnit"),
		ItemBillingBlockReason:  getString(raw, "itemBillingBlockReason"),
		LastChangeDate:          lastChangeDate,
		Plant:                   getString(raw, "plant"),
		ReferenceSdDocument:     getString(raw, "referenceSdDocument"),
		ReferenceSdDocumentItem: getString(raw, "referenceSdDocumentItem"),
		StorageLocation:         getString(raw, "storageLocation"),
	}
	return m, nil
}
