package loader

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/sudarsh1010/graph-rag/internal/model"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type ProductsLoader struct{ baseLoader }

func NewProductsLoader() *ProductsLoader {
	return &ProductsLoader{baseLoader: newBaseLoader("products", "products", "products")}
}

func (l *ProductsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *ProductsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	crossPlantDate, _ := getTimePtr(raw, "crossPlantStatusValidityDate")
	creationDate, _ := getTime(raw, "creationDate")
	lastChangeDate, _ := getTime(raw, "lastChangeDate")
	lastChangeDateTime, _ := getTime(raw, "lastChangeDateTime")

	isMarkedForDeletion := getBool(raw, "isMarkedForDeletion")

	m := &model.Product{
		Product:                      getString(raw, "product"),
		ProductType:                  getString(raw, "productType"),
		CrossPlantStatus:             getString(raw, "crossPlantStatus"),
		CrossPlantStatusValidityDate: crossPlantDate,
		CreationDate:                 creationDate,
		CreatedByUser:                getString(raw, "createdByUser"),
		LastChangeDate:               lastChangeDate,
		LastChangeDateTime:           lastChangeDateTime,
		IsMarkedForDeletion:          &isMarkedForDeletion,
		ProductOldId:                 getString(raw, "productOldId"),
		GrossWeight:                  mustDecimal(raw, "grossWeight"),
		WeightUnit:                   getString(raw, "weightUnit"),
		NetWeight:                    mustDecimal(raw, "netWeight"),
		ProductGroup:                 getString(raw, "productGroup"),
		BaseUnit:                     getString(raw, "baseUnit"),
		Division:                     getString(raw, "division"),
		IndustrySector:               getString(raw, "industrySector"),
	}
	if err := ValidateRequired(m.Product, "product"); err != nil {
		return nil, err
	}
	return m, nil
}

func mustDecimal(raw map[string]interface{}, key string) decimal.Decimal {
	d, _ := getDecimal(raw, key)
	return d
}
