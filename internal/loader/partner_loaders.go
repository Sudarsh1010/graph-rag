package loader

import (
	"context"

	"github.com/sudarsh1010/graph-rag/internal/model"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type BusinessPartnersLoader struct{ baseLoader }

func NewBusinessPartnersLoader() *BusinessPartnersLoader {
	return &BusinessPartnersLoader{baseLoader: newBaseLoader("business_partners", "business_partners", "business_partners")}
}

func (l *BusinessPartnersLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *BusinessPartnersLoader) convert(raw map[string]interface{}) (interface{}, error) {
	creationDate, _ := getTime(raw, "creationDate")
	lastChangeDate, _ := getTime(raw, "lastChangeDate")

	m := &model.BusinessPartner{
		BusinessPartner:          getString(raw, "businessPartner"),
		Customer:                 getString(raw, "customer"),
		BusinessPartnerCategory:  getString(raw, "businessPartnerCategory"),
		BusinessPartnerFullName:  getString(raw, "businessPartnerFullName"),
		BusinessPartnerGrouping:  getString(raw, "businessPartnerGrouping"),
		BusinessPartnerName:      getString(raw, "businessPartnerName"),
		CorrespondenceLanguage:   getStringPtr(raw, "correspondenceLanguage"),
		CreatedByUser:            getString(raw, "createdByUser"),
		CreationDate:             creationDate,
		CreationTime:             mustTimeOfDay(raw, "creationTime"),
		FirstName:                getStringPtr(raw, "firstName"),
		FormOfAddress:            getString(raw, "formOfAddress"),
		Industry:                 getStringPtr(raw, "industry"),
		LastChangeDate:           lastChangeDate,
		LastName:                 getStringPtr(raw, "lastName"),
		OrganizationBPName1:      getString(raw, "organizationBpName1"),
		OrganizationBPName2:      getStringPtr(raw, "organizationBpName2"),
		BusinessPartnerIsBlocked: getBool(raw, "businessPartnerIsBlocked"),
		IsMarkedForArchiving:     getBool(raw, "isMarkedForArchiving"),
	}
	if err := ValidateRequired(m.BusinessPartner, "businessPartner"); err != nil {
		return nil, err
	}
	return m, nil
}

func mustTimeOfDay(raw map[string]interface{}, key string) model.TimeOfDay {
	tod, _ := getTimeOfDay(raw, key)
	return tod
}

type BusinessPartnerAddressesLoader struct{ baseLoader }

func NewBusinessPartnerAddressesLoader() *BusinessPartnerAddressesLoader {
	return &BusinessPartnerAddressesLoader{baseLoader: newBaseLoader("business_partner_addresses", "business_partner_addresses", "business_partner_addresses")}
}

func (l *BusinessPartnerAddressesLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *BusinessPartnerAddressesLoader) convert(raw map[string]interface{}) (interface{}, error) {
	validityStart, _ := getTime(raw, "validityStartDate")
	validityEnd, _ := getTime(raw, "validityEndDate")

	m := &model.BusinessPartnerAddress{
		BusinessPartner:        getString(raw, "businessPartner"),
		AddressID:              getString(raw, "addressId"),
		ValidityStartDate:      validityStart,
		ValidityEndDate:        validityEnd,
		AddressUUID:            getString(raw, "addressUuid"),
		AddressTimeZone:        getString(raw, "addressTimeZone"),
		CityName:               getStringPtr(raw, "cityName"),
		Country:                getString(raw, "country"),
		POBox:                  getStringPtr(raw, "poBox"),
		POBoxDeviatingCityName: getStringPtr(raw, "poBoxDeviatingCityName"),
		POBoxDeviatingCountry:  getStringPtr(raw, "poBoxDeviatingCountry"),
		POBoxDeviatingRegion:   getStringPtr(raw, "poBoxDeviatingRegion"),
		POBoxIsWithoutNumber:   getBool(raw, "poBoxIsWithoutNumber"),
		POBoxLobbyName:         getStringPtr(raw, "poBoxLobbyName"),
		POBoxPostalCode:        getStringPtr(raw, "poBoxPostalCode"),
		PostalCode:             getStringPtr(raw, "postalCode"),
		Region:                 getStringPtr(raw, "region"),
		StreetName:             getStringPtr(raw, "streetName"),
		TaxJurisdiction:        getStringPtr(raw, "taxJurisdiction"),
		TransportZone:          getStringPtr(raw, "transportZone"),
	}
	return m, nil
}

type CustomerCompanyAssignmentsLoader struct{ baseLoader }

func NewCustomerCompanyAssignmentsLoader() *CustomerCompanyAssignmentsLoader {
	return &CustomerCompanyAssignmentsLoader{baseLoader: newBaseLoader("customer_company_assignments", "customer_company_assignments", "customer_company_assignments")}
}

func (l *CustomerCompanyAssignmentsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *CustomerCompanyAssignmentsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	m := &model.CustomerCompanyAssignment{
		Customer:                    getString(raw, "customer"),
		CompanyCode:                 getString(raw, "companyCode"),
		AccountingClerk:             getStringPtr(raw, "accountingClerk"),
		AccountingClerkFaxNumber:    getStringPtr(raw, "accountingClerkFaxNumber"),
		AccountingClerkEmailAddress: getStringPtr(raw, "accountingClerkInternetAddress"),
		AccountingClerkPhoneNumber:  getStringPtr(raw, "accountingClerkPhoneNumber"),
		AlternativePayerAccount:     getStringPtr(raw, "alternativePayerAccount"),
		PaymentBlockingReason:       getStringPtr(raw, "paymentBlockingReason"),
		PaymentMethodsList:          getStringPtr(raw, "paymentMethodsList"),
		PaymentTerms:                getStringPtr(raw, "paymentTerms"),
		ReconciliationAccount:       getString(raw, "reconciliationAccount"),
		DeletionIndicator:           getBool(raw, "deletionIndicator"),
		CustomerAccountGroup:        getString(raw, "customerAccountGroup"),
	}
	return m, nil
}

type CustomerSalesAreaAssignmentsLoader struct{ baseLoader }

func NewCustomerSalesAreaAssignmentsLoader() *CustomerSalesAreaAssignmentsLoader {
	return &CustomerSalesAreaAssignmentsLoader{baseLoader: newBaseLoader("customer_sales_area_assignments", "customer_sales_area_assignments", "customer_sales_area_assignments")}
}

func (l *CustomerSalesAreaAssignmentsLoader) Load(ctx context.Context, db *bun.DB, datasetPath string, logger *zap.Logger) (*LoadResult, error) {
	return l.loadFiles(ctx, db, datasetPath, logger, l.convert, 100)
}

func (l *CustomerSalesAreaAssignmentsLoader) convert(raw map[string]interface{}) (interface{}, error) {
	m := &model.CustomerSalesAreaAssignment{
		Customer:                    getString(raw, "customer"),
		SalesOrganization:           getString(raw, "salesOrganization"),
		DistributionChannel:         getString(raw, "distributionChannel"),
		Division:                    getString(raw, "division"),
		BillingIsBlockedForCustomer: getStringPtr(raw, "billingIsBlockedForCustomer"),
		CompleteDeliveryIsDefined:   getBool(raw, "completeDeliveryIsDefined"),
		CreditControlArea:           getStringPtr(raw, "creditControlArea"),
		Currency:                    getString(raw, "currency"),
		CustomerPaymentTerms:        getString(raw, "customerPaymentTerms"),
		DeliveryPriority:            getString(raw, "deliveryPriority"),
		IncotermsClassification:     getString(raw, "incotermsClassification"),
		IncotermsLocation1:          getString(raw, "incotermsLocation1"),
		SalesGroup:                  getStringPtr(raw, "salesGroup"),
		SalesOffice:                 getStringPtr(raw, "salesOffice"),
		ShippingCondition:           getString(raw, "shippingCondition"),
		SlsUnlmtdOvrdelivIsAllwd:    getBool(raw, "slsUnlmtdOvrdelivIsAllwd"),
		SupplyingPlant:              getStringPtr(raw, "supplyingPlant"),
		SalesDistrict:               getStringPtr(raw, "salesDistrict"),
		ExchangeRateType:            getStringPtr(raw, "exchangeRateType"),
	}
	return m, nil
}
