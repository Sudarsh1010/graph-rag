package loader

import (
	"context"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type Coordinator struct {
	loaders []EntityLoader
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		loaders: []EntityLoader{
			NewPlantsLoader(),
			NewProductsLoader(),
			NewBusinessPartnersLoader(),
			NewProductDescriptionsLoader(),
			NewProductPlantsLoader(),
			NewProductStorageLocationsLoader(),
			NewBusinessPartnerAddressesLoader(),
			NewCustomerCompanyAssignmentsLoader(),
			NewCustomerSalesAreaAssignmentsLoader(),
			NewSalesOrderHeadersLoader(),
			NewOutboundDeliveryHeadersLoader(),
			NewSalesOrderItemsLoader(),
			NewOutboundDeliveryItemsLoader(),
			NewSalesOrderScheduleLinesLoader(),
			NewBillingDocumentHeadersLoader(),
			NewBillingDocumentItemsLoader(),
			NewBillingDocumentCancellationsLoader(),
			NewJournalEntryItemsARLoader(),
			NewPaymentsARLoader(),
		},
	}
}

type LoadConfig struct {
	DatasetPath string
	DB          *bun.DB
	Logger      *zap.Logger
}

func (c *Coordinator) LoadAll(ctx context.Context, cfg LoadConfig) (*LoadResult, error) {
	logger := cfg.Logger
	combined := &LoadResult{Entity: "all"}

	for _, loader := range c.loaders {
		result, err := loader.Load(ctx, cfg.DB, cfg.DatasetPath, logger)
		if err != nil {
			logger.Warn("loader failed, continuing",
				zap.String("entity", loader.EntityName()),
				zap.Error(err),
			)
		}
		if result != nil {
			combined.Merge(result)
		}
	}

	logger.Info("completed loading all entities",
		zap.Int("total", combined.Total),
		zap.Int("inserted", combined.Inserted),
		zap.Int("skipped", combined.Skipped),
		zap.Int("errors", len(combined.Errors)),
	)

	return combined, nil
}
