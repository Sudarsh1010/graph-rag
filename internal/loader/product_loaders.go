package loader

import (
	"context"

	"github.com/sudarsh1010/graph-rag/internal/model"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type ProductDescriptionsLoader struct{ baseLoader }

func NewProductDescriptionsLoader() *ProductDescriptionsLoader {
	return &ProductDescriptionsLoader{baseLoader: newBaseLoader("product_descriptions", "product_descriptions", "product_descriptions")}
}

func (l *ProductDescriptionsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *ProductDescriptionsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	m := &model.ProductDescription{
		Product:            getString(raw, "product"),
		Language:           getString(raw, "language"),
		ProductDescription: getString(raw, "productDescription"),
	}
	return m, nil
}

type ProductPlantsLoader struct{ baseLoader }

func NewProductPlantsLoader() *ProductPlantsLoader {
	return &ProductPlantsLoader{baseLoader: newBaseLoader("product_plants", "product_plants", "product_plants")}
}

func (l *ProductPlantsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 500)
}

func (l *ProductPlantsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	m := &model.ProductPlant{
		Product:                    getString(raw, "product"),
		Plant:                      getString(raw, "plant"),
		CountryOfOrigin:            getString(raw, "countryOfOrigin"),
		RegionOfOrigin:             getString(raw, "regionOfOrigin"),
		ProductionInvtryManagedLoc: getString(raw, "productionInvtryManagedLoc"),
		AvailabilityCheckType:      getString(raw, "availabilityCheckType"),
		FiscalYearVariant:          getString(raw, "fiscalYearVariant"),
		ProfitCenter:               getString(raw, "profitCenter"),
		MrpType:                    getString(raw, "mrpType"),
	}
	return m, nil
}

type ProductStorageLocationsLoader struct{ baseLoader }

func NewProductStorageLocationsLoader() *ProductStorageLocationsLoader {
	return &ProductStorageLocationsLoader{baseLoader: newBaseLoader("product_storage_locations", "product_storage_locations", "product_storage_locations")}
}

func (l *ProductStorageLocationsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 500)
}

func (l *ProductStorageLocationsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	lastPosted, _ := getTimePtr(raw, "dateOfLastPostedCntUnRstrcdStk")

	m := &model.ProductStorageLocation{
		Product:                        getString(raw, "product"),
		Plant:                          getString(raw, "plant"),
		StorageLocation:                getString(raw, "storageLocation"),
		PhysicalInventoryBlockInd:      getString(raw, "physicalInventoryBlockInd"),
		DateOfLastPostedCntUnRstrcdStk: lastPosted,
	}
	return m, nil
}
