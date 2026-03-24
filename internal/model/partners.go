package model

import (
	"time"

	"github.com/uptrace/bun"
)

// BusinessPartner represents customer/partner master data.
// Source: business_partners JSONL files
type BusinessPartner struct {
	bun.BaseModel `bun:"table:business_partners,alias:bp"`

	BusinessPartner          string    `json:"businessPartner" bun:",pk"`
	Customer                 string    `json:"customer"`
	BusinessPartnerCategory  string    `json:"businessPartnerCategory"`
	BusinessPartnerFullName  string    `json:"businessPartnerFullName"`
	BusinessPartnerGrouping  string    `json:"businessPartnerGrouping"`
	BusinessPartnerName      string    `json:"businessPartnerName"`
	CorrespondenceLanguage   *string   `json:"correspondenceLanguage,omitempty" bun:",nullzero"`
	CreatedByUser            string    `json:"createdByUser"`
	CreationDate             time.Time `json:"creationDate"`
	CreationTime             TimeOfDay `json:"creationTime"`
	FirstName                *string   `json:"firstName,omitempty" bun:",nullzero"`
	FormOfAddress            string    `json:"formOfAddress"`
	Industry                 *string   `json:"industry,omitempty" bun:",nullzero"`
	LastChangeDate           time.Time `json:"lastChangeDate"`
	LastName                 *string   `json:"lastName,omitempty" bun:",nullzero"`
	OrganizationBPName1      string    `json:"organizationBpName1"`
	OrganizationBPName2      *string   `json:"organizationBpName2,omitempty" bun:",nullzero"`
	BusinessPartnerIsBlocked bool      `json:"businessPartnerIsBlocked"`
	IsMarkedForArchiving     bool      `json:"isMarkedForArchiving"`
}

// BusinessPartnerAddress represents address data for business partners.
// Source: business_partner_addresses JSONL files
type BusinessPartnerAddress struct {
	bun.BaseModel `bun:"table:business_partner_addresses,alias:bpa"`

	BusinessPartner        string    `json:"businessPartner" bun:",pk"`
	AddressID              string    `json:"addressId" bun:",pk"`
	ValidityStartDate      time.Time `json:"validityStartDate"`
	ValidityEndDate        time.Time `json:"validityEndDate"`
	AddressUUID            string    `json:"addressUuid"`
	AddressTimeZone        string    `json:"addressTimeZone"`
	CityName               *string   `json:"cityName,omitempty" bun:",nullzero"`
	Country                string    `json:"country"`
	POBox                  *string   `json:"poBox,omitempty" bun:",nullzero"`
	POBoxDeviatingCityName *string   `json:"poBoxDeviatingCityName,omitempty" bun:",nullzero"`
	POBoxDeviatingCountry  *string   `json:"poBoxDeviatingCountry,omitempty" bun:",nullzero"`
	POBoxDeviatingRegion   *string   `json:"poBoxDeviatingRegion,omitempty" bun:",nullzero"`
	POBoxIsWithoutNumber   bool      `json:"poBoxIsWithoutNumber"`
	POBoxLobbyName         *string   `json:"poBoxLobbyName,omitempty" bun:",nullzero"`
	POBoxPostalCode        *string   `json:"poBoxPostalCode,omitempty" bun:",nullzero"`
	PostalCode             *string   `json:"postalCode,omitempty" bun:",nullzero"`
	Region                 *string   `json:"region,omitempty" bun:",nullzero"`
	StreetName             *string   `json:"streetName,omitempty" bun:",nullzero"`
	TaxJurisdiction        *string   `json:"taxJurisdiction,omitempty" bun:",nullzero"`
	TransportZone          *string   `json:"transportZone,omitempty" bun:",nullzero"`
}

// CustomerCompanyAssignment represents company code assignments for customers.
// Source: customer_company_assignments JSONL files
type CustomerCompanyAssignment struct {
	bun.BaseModel `bun:"table:customer_company_assignments,alias:cca"`

	Customer                    string  `json:"customer" bun:",pk"`
	CompanyCode                 string  `json:"companyCode" bun:",pk"`
	AccountingClerk             *string `json:"accountingClerk,omitempty" bun:",nullzero"`
	AccountingClerkFaxNumber    *string `json:"accountingClerkFaxNumber,omitempty" bun:",nullzero"`
	AccountingClerkEmailAddress *string `json:"accountingClerkInternetAddress,omitempty" bun:",nullzero"`
	AccountingClerkPhoneNumber  *string `json:"accountingClerkPhoneNumber,omitempty" bun:",nullzero"`
	AlternativePayerAccount     *string `json:"alternativePayerAccount,omitempty" bun:",nullzero"`
	PaymentBlockingReason       *string `json:"paymentBlockingReason,omitempty" bun:",nullzero"`
	PaymentMethodsList          *string `json:"paymentMethodsList,omitempty" bun:",nullzero"`
	PaymentTerms                *string `json:"paymentTerms,omitempty" bun:",nullzero"`
	ReconciliationAccount       string  `json:"reconciliationAccount"`
	DeletionIndicator           bool    `json:"deletionIndicator"`
	CustomerAccountGroup        string  `json:"customerAccountGroup"`
}

// CustomerSalesAreaAssignment represents sales area data for customers.
// Source: customer_sales_area_assignments JSONL files
type CustomerSalesAreaAssignment struct {
	bun.BaseModel `bun:"table:customer_sales_area_assignments,alias:csaa"`

	Customer                    string  `json:"customer" bun:",pk"`
	SalesOrganization           string  `json:"salesOrganization" bun:",pk"`
	DistributionChannel         string  `json:"distributionChannel" bun:",pk"`
	Division                    string  `json:"division" bun:",pk"`
	BillingIsBlockedForCustomer *string `json:"billingIsBlockedForCustomer,omitempty" bun:",nullzero"`
	CompleteDeliveryIsDefined   bool    `json:"completeDeliveryIsDefined"`
	CreditControlArea           *string `json:"creditControlArea,omitempty" bun:",nullzero"`
	Currency                    string  `json:"currency"`
	CustomerPaymentTerms        string  `json:"customerPaymentTerms"`
	DeliveryPriority            string  `json:"deliveryPriority"`
	IncotermsClassification     string  `json:"incotermsClassification"`
	IncotermsLocation1          string  `json:"incotermsLocation1"`
	SalesGroup                  *string `json:"salesGroup,omitempty" bun:",nullzero"`
	SalesOffice                 *string `json:"salesOffice,omitempty" bun:",nullzero"`
	ShippingCondition           string  `json:"shippingCondition"`
	SlsUnlmtdOvrdelivIsAllwd    bool    `json:"slsUnlmtdOvrdelivIsAllwd"`
	SupplyingPlant              *string `json:"supplyingPlant,omitempty" bun:",nullzero"`
	SalesDistrict               *string `json:"salesDistrict,omitempty" bun:",nullzero"`
	ExchangeRateType            *string `json:"exchangeRateType,omitempty" bun:",nullzero"`
}
