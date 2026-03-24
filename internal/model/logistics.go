package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type OutboundDeliveryHeader struct {
	ActualGoodsMovementDate      *time.Time `bun:",nullzero" json:"actualGoodsMovementDate,omitempty"`
	ActualGoodsMovementTime      TimeOfDay  `json:"actualGoodsMovementTime,omitempty"`
	CreationDate                 string     `json:"creationDate,omitempty"`
	CreationTime                 TimeOfDay  `json:"creationTime,omitempty"`
	DeliveryBlockReason          string     `json:"deliveryBlockReason,omitempty"`
	DeliveryDocument             string     `bun:",pk" json:"deliveryDocument,omitempty"`
	HdrGeneralIncompletionStatus string     `json:"hdrGeneralIncompletionStatus,omitempty"`
	HeaderBillingBlockReason     string     `json:"headerBillingBlockReason,omitempty"`
	LastChangeDate               *time.Time `bun:",nullzero" json:"lastChangeDate,omitempty"`
	OverallGoodsMovementStatus   string     `json:"overallGoodsMovementStatus,omitempty"`
	OverallPickingStatus         string     `json:"overallPickingStatus,omitempty"`
	OverallProofOfDeliveryStatus string     `json:"overallProofOfDeliveryStatus,omitempty"`
	ShippingPoint                string     `json:"shippingPoint,omitempty"`
}

type OutboundDeliveryItem struct {
	ActualDeliveryQuantity  decimal.Decimal `bun:",nullzero" json:"actualDeliveryQuantity,omitempty"`
	Batch                   string          `json:"batch,omitempty"`
	DeliveryDocument        string          `bun:",pk" json:"deliveryDocument,omitempty"`
	DeliveryDocumentItem    string          `bun:",pk" json:"deliveryDocumentItem,omitempty"`
	DeliveryQuantityUnit    string          `json:"deliveryQuantityUnit,omitempty"`
	ItemBillingBlockReason  string          `json:"itemBillingBlockReason,omitempty"`
	LastChangeDate          *time.Time      `bun:",nullzero" json:"lastChangeDate,omitempty"`
	Plant                   string          `json:"plant,omitempty"`
	ReferenceSdDocument     string          `json:"referenceSdDocument,omitempty"`
	ReferenceSdDocumentItem string          `json:"referenceSdDocumentItem,omitempty"`
	StorageLocation         string          `json:"storageLocation,omitempty"`
}

type Plant struct {
	Plant                         string `bun:",pk" json:"plant,omitempty"`
	PlantName                     string `json:"plantName,omitempty"`
	ValuationArea                 string `json:"valuationArea,omitempty"`
	PlantCustomer                 string `json:"plantCustomer,omitempty"`
	PlantSupplier                 string `json:"plantSupplier,omitempty"`
	FactoryCalendar               string `json:"factoryCalendar,omitempty"`
	DefaultPurchasingOrganization string `json:"defaultPurchasingOrganization,omitempty"`
	SalesOrganization             string `json:"salesOrganization,omitempty"`
	AddressId                     string `json:"addressId,omitempty"`
	PlantCategory                 string `json:"plantCategory,omitempty"`
	DistributionChannel           string `json:"distributionChannel,omitempty"`
	Division                      string `json:"division,omitempty"`
	Language                      string `json:"language,omitempty"`
	IsMarkedForArchiving          bool   `json:"isMarkedForArchiving,omitempty"`
}
