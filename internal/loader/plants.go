package loader

import (
	"context"

	"github.com/sudarsh1010/graph-rag/internal/model"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type PlantsLoader struct{ baseLoader }

func NewPlantsLoader() *PlantsLoader {
	return &PlantsLoader{baseLoader: newBaseLoader("plants", "plants", "plants")}
}

func (l *PlantsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *PlantsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	m := &model.Plant{
		Plant:                         getString(raw, "plant"),
		PlantName:                     getString(raw, "plantName"),
		ValuationArea:                 getString(raw, "valuationArea"),
		PlantCustomer:                 getString(raw, "plantCustomer"),
		PlantSupplier:                 getString(raw, "plantSupplier"),
		FactoryCalendar:               getString(raw, "factoryCalendar"),
		DefaultPurchasingOrganization: getString(raw, "defaultPurchasingOrganization"),
		SalesOrganization:             getString(raw, "salesOrganization"),
		AddressId:                     getString(raw, "addressId"),
		PlantCategory:                 getString(raw, "plantCategory"),
		DistributionChannel:           getString(raw, "distributionChannel"),
		Division:                      getString(raw, "division"),
		Language:                      getString(raw, "language"),
		IsMarkedForArchiving:          getBool(raw, "isMarkedForArchiving"),
	}
	if err := ValidateRequired(m.Plant, "plant"); err != nil {
		return nil, err
	}
	return m, nil
}
