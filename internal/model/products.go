package model

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	Product                      string          `bun:"type:varchar(40),pk,nullzero,notnull"`
	ProductType                  string          `bun:"type:varchar(10),nullzero"`
	CrossPlantStatus             string          `bun:"type:varchar(1),nullzero"`
	CrossPlantStatusValidityDate *time.Time      `bun:"type:timestamp,nullzero" json:",omitempty"`
	CreationDate                 time.Time       `bun:"type:timestamp,nullzero,notnull"`
	CreatedByUser                string          `bun:"type:varchar(10),nullzero"`
	LastChangeDate               time.Time       `bun:"type:timestamp,nullzero,notnull"`
	LastChangeDateTime           time.Time       `bun:"type:timestamp,nullzero,notnull"`
	IsMarkedForDeletion          *bool           `bun:"type:boolean,nullzero" json:",omitempty"`
	ProductOldId                 string          `bun:"type:varchar(40),nullzero"`
	GrossWeight                  decimal.Decimal `bun:"type:decimal(18,4),nullzero"`
	WeightUnit                   string          `bun:"type:varchar(3),nullzero"`
	NetWeight                    decimal.Decimal `bun:"type:decimal(18,4),nullzero"`
	ProductGroup                 string          `bun:"type:varchar(10),nullzero"`
	BaseUnit                     string          `bun:"type:varchar(4),nullzero"`
	Division                     string          `bun:"type:varchar(2),nullzero"`
	IndustrySector               string          `bun:"type:varchar(1),nullzero"`
}

type ProductPlant struct {
	bun.BaseModel              `bun:"table:product_plants,alias:pp"`
	Product                    string `bun:"type:varchar(40),pk,nullzero,notnull"`
	Plant                      string `bun:"type:varchar(4),pk,nullzero,notnull"`
	CountryOfOrigin            string `bun:"type:varchar(3),nullzero"`
	RegionOfOrigin             string `bun:"type:varchar(4),nullzero"`
	ProductionInvtryManagedLoc string `bun:"type:varchar(40),nullzero"`
	AvailabilityCheckType      string `bun:"type:varchar(2),nullzero"`
	FiscalYearVariant          string `bun:"type:varchar(2),nullzero"`
	ProfitCenter               string `bun:"type:varchar(10),nullzero"`
	MrpType                    string `bun:"type:varchar(2),nullzero"`
}

type ProductStorageLocation struct {
	bun.BaseModel                  `bun:"table:product_storage_locations,alias:psl"`
	Product                        string     `bun:"type:varchar(40),pk,nullzero,notnull"`
	Plant                          string     `bun:"type:varchar(4),pk,nullzero,notnull"`
	StorageLocation                string     `bun:"type:varchar(10),pk,nullzero,notnull"`
	PhysicalInventoryBlockInd      string     `bun:"type:varchar(1),nullzero"`
	DateOfLastPostedCntUnRstrcdStk *time.Time `bun:"type:timestamp,nullzero" json:",omitempty"`
}

type ProductDescription struct {
	bun.BaseModel      `bun:"table:product_descriptions,alias:pd"`
	Product            string `bun:"type:varchar(40),pk,nullzero,notnull"`
	Language           string `bun:"type:varchar(2),pk,nullzero,notnull"`
	ProductDescription string `bun:"type:varchar(300),nullzero"`
}
