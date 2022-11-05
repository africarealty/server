package domain

import (
	"context"
	"github.com/africarealty/server/src/kit"
	"time"
)

const (
	AdStatusDraft  = "draft"
	AdStatusOpen   = "open"
	AdStatusClosed = "closed"

	AdSubStatusDraft                  = "draft"
	AdSubStatusDraftReview            = "review"
	AdSubStatusOpenActive             = "active"
	AdSubStatusOpenSuspended          = "suspended"
	AdSubStatusClosedCancelledByOwner = "cancelled-by-owner"
	AdSubStatusClosedCancelledByAdmin = "cancelled-by-admin"
	AdSubStatusClosedExpired          = "expired"
	AdSubStatusClosedResolved         = "resolved"

	AdTypeRealty        = "realty"
	AdSubTypeRealtyRent = "rent"
	AdSubTypeRealtySale = "sale"
)

type AdAdvertiserType uint

const (
	AdAdvertiserOwner AdAdvertiserType = iota
	AdAdvertiserAgent
)

type AdRealtyDealTermType uint

const (
	AdRealtyDealTermLong AdRealtyDealTermType = iota
	AdRealtyDealTermMedium
	AdRealtyDealTermShort
)

type AdRealtyType uint

const (
	AdRealtyFlat AdRealtyType = iota
	AdRealtyApartment
	AdRealtyRoom
	AdRealtyBed
	AdRealtyHouse
)

type AdRealtyStructureType uint

const (
	AdRealtyStructureStudio AdRealtyStructureType = iota
	AdRealtyStructureOne
	AdRealtyStructureOneAndHalf
	AdRealtyStructureTwo
	AdRealtyStructureTwoAndHalf
	AdRealtyStructureThree
	AdRealtyStructureThreeAndHalf
	AdRealtyStructureFour
	AdRealtyStructureFourAndHalf
	AdRealtyStructureFivePlus
)

type AdRealtyNumberOfBedroomsType uint

const (
	AdRealtyNumberOfBedroomsOne AdRealtyNumberOfBedroomsType = iota + 1
	AdRealtyNumberOfBedroomsTwo
	AdRealtyNumberOfBedroomsThree
	AdRealtyNumberOfBedroomsFour
	AdRealtyNumberOfBedroomsFive
	AdRealtyNumberOfBedroomsSixAndMore
)

type AdRealtyNumberOfBathroomsType uint

const (
	AdRealtyNumberOfBathroomsOne AdRealtyNumberOfBathroomsType = iota + 1
	AdRealtyNumberOfBathroomsTwo
	AdRealtyNumberOfBathroomsThree
	AdRealtyNumberOfBathroomsFourAndMore
)

type AdRealtyParkingType uint

const (
	AdRealtyParkingNo AdRealtyParkingType = iota
	AdRealtyParkingGarage
	AdRealtyParkingFreeZone
	AdRealtyParkingPrivateParking
	AdRealtyParkingBuildingParking
)

type AdRealtyRenovationType uint

const (
	AdRealtyRenovationNo AdRealtyRenovationType = iota
	AdRealtyRenovationByDesigner
	AdRealtyRenovationEuro
	AdRealtyRenovationRegular
)

type AdRealtyWindowViewType = uint

const (
	AdRealtyWindowViewYard AdRealtyWindowViewType = iota
	AdRealtyWindowViewSea
	AdRealtyWindowViewStreet
)

type NearByType = uint

const (
	NearByLessThan5 NearByType = iota
	NearByFrom5to10
	NearByMoreThan10
)

type FurnishedType = uint

const (
	UnFurnished FurnishedType = iota
	SemiFurnished
	Furnished
)

// Geo position
type Geo struct {
	Lat  float32 `json:"lat"`  // Lat latitude
	Long float32 `json:"long"` // Long longitude
}

// Address of realty
type Address struct {
	Full        string  `json:"full"`                 // Full address
	Short       string  `json:"short"`                // Short address
	PostCode    *uint   `json:"postCode,omitempty"`   // PostCode post code
	Country     string  `json:"country"`              // Country name
	CountryCode string  `json:"countryCode"`          // CountryCode country code (country code-table)
	Region      string  `json:"region,omitempty"`     // Region name
	RegionCode  string  `json:"regionCode,omitempty"` // RegionCode region code (region code-table)
	City        string  `json:"city,omitempty"`       // City name
	CityCode    string  `json:"cityCode,omitempty"`   // CityCode city code
	Street      string  `json:"street,omitempty"`     // Street name
	Building    string  `json:"building,omitempty"`   // Building number
	Apartment   string  `json:"apartment,omitempty"`  // Apartment number
	Floor       *uint32 `json:"floor,omitempty"`      // Floor number
	Entrance    *uint32 `json:"entrance,omitempty"`   // Entrance number
	Geo         *Geo    `json:"geo,omitempty"`        // Geo position
}

// ParticipantPhone phone of participant
type ParticipantPhone struct {
	Number           string `json:"number"`                     // Number phone number
	Sms              *bool  `json:"sms,omitempty"`              // Sms if sms available
	Calls            *bool  `json:"calls,omitempty"`            // Calls if calls available
	WhatsUp          *bool  `json:"whatsUp,omitempty"`          // WhatsUp if WhatsUp available
	Viber            *bool  `json:"viber,omitempty"`            // Viber if Viber available
	Telegram         *bool  `json:"telegram,omitempty"`         // Telegram if Telegram available
	AvailabilityTime string `json:"availabilityTime,omitempty"` // AvailabilityTime when available for contacts
}

// Participant details
type Participant struct {
	Name   string            `json:"name,omitempty"`   // Name participant
	UserId string            `json:"userId,omitempty"` // UserId user ID
	Email  string            `json:"email,omitempty"`  // Email to contact
	Phone  *ParticipantPhone `json:"phone,omitempty"`  // Phone to contact
}

// Participants list
type Participants struct {
	Advertiser *Participant `json:"advertiser,omitempty"` // Advertiser manages ads
	Reviewer   *Participant `json:"reviewer,omitempty"`   // Reviewer reviews ads
	Owner      *Participant `json:"owner,omitempty"`      // Owner owns realty
	Agent      *Participant `json:"agent,omitempty"`      // Agent deal agent
}

// Description ads description
type Description struct {
	Title string `json:"title,omitempty"` // Title contains ads summary
	Short string `json:"short,omitempty"` // Short describes ads shortly
	Full  string `json:"full,omitempty"`  // Full description
}

// RentPeriod specifies rent period
type RentPeriod struct {
	Unit TimeUnitType `json:"unit,omitempty"` // Unit period unit
	Min  *uint32      `json:"min,omitempty"`  // Min minimum period
	Max  *uint32      `json:"max,omitempty"`  // Max maximum period
}

type RegularPayment struct {
	Unit   PeriodUnitType `json:"unit,omitempty"`   // Unit period unit
	Cur    Currency       `json:"cur,omitempty"`    // Cur currency
	Amount float64        `json:"amount,omitempty"` // Amount per unit period
}

// Expense data
type Expense struct {
	PerConsumption bool            `json:"perConsumption"`       // PerConsumption if true, expense is paid by consumption
	RegularPayment *RegularPayment `json:"regularPay,omitempty"` // RegularPayment specified if expense has regularly paid
}

// RentExpenses rent expenses
type RentExpenses struct {
	Electricity         *Expense `json:"electricity,omitempty"`         // Electricity expense
	Heating             *Expense `json:"heating,omitempty"`             // Heating expense
	Water               *Expense `json:"water,omitempty"`               // Water expense
	BuildingMaintenance *Expense `json:"buildingMaintenance,omitempty"` // BuildingMaintenance expense
	Internet            *Expense `json:"internet,omitempty"`            // Internet expense
	Total               *Expense `json:"total,omitempty"`               // Total expense
}

// RentDeposit rent deposit
type RentDeposit struct {
	Deposit *float64 `json:"deposit,omitempty"` // Deposit deposit amount
	Cur     Currency `json:"cur,omitempty"`     // Cur deposit currency
}

// RentTerms terms of rent
type RentTerms struct {
	Price       *RegularPayment `json:"price,omitempty"`       // Price rent price
	Deposit     *RentDeposit    `json:"deposit,omitempty"`     // Deposit required for rent
	Period      *RentPeriod     `json:"period,omitempty"`      // Period rent period
	Expenses    *RentExpenses   `json:"expenses,omitempty"`    // Expenses rent expenses
	PetsAllowed *bool           `json:"petsAllowed,omitempty"` // PetsAllowed pets allowed
	KidsAllowed *bool           `json:"kidsAllowed,omitempty"` // KidsAllowed kids allowed
}

// NearBy specifies what places are around the building
type NearBy struct {
	SeaSide    *NearByType `json:"sea,omitempty"`        // Distance to SeaSide in meters
	CityCenter *NearByType `json:"cityCenter,omitempty"` // Distance to CityCenter in meters
	Mall       *NearByType `json:"mall,omitempty"`       // Distance to Mall in meters
	Park       *NearByType `json:"park,omitempty"`       // Distance to Park in meters
}

type Parking struct {
	Type AdRealtyParkingType
}

// Location describes building location
type Location struct {
	Address *Address `json:"address,omitempty"` // Building Address
	NearBy  *NearBy  `json:"nearBy,omitempty"`  // What objects are NearBy the building
	Parking *Parking `json:"parking,omitempty"` // Parking details
}

// Building specifies building details
type Building struct {
	Floors         *uint32 `json:"floors,omitempty"`         // Number of Floors
	Year           *uint32 `json:"year,omitempty"`           // Year of building
	LastRenovation *uint32 `json:"lastRenovation,omitempty"` // LastRenovation year
	Elevator       *bool   `json:"elevator,omitempty"`       // If Elevator exists
	Security       *bool   `json:"security,omitempty"`       // If Security service exists
	Cameras        *bool   `json:"cameras,omitempty"`        // If Cameras exist
	SwimmingPool   *bool   `json:"swimmingPool,omitempty"`   // If private SwimmingPool exist
	Intercom       *bool   `json:"intercom,omitempty"`       // If Intercom exist
}

// Area of the realty object
type Area struct {
	Total    float32  `json:"total,omitempty"`    // Total area in sq meters
	Living   *float32 `json:"living,omitempty"`   // Area of Living room in sq meters
	Kitchen  *float32 `json:"kitchen,omitempty"`  // Area of Kitchen in sq meters
	BedRooms *float32 `json:"bedRooms,omitempty"` // Area of BedRooms in sq meters
	Balcony  *float32 `json:"balcony,omitempty"`  // Area of Balcony in sq meters
}

// Equipment in realty
type Equipment struct {
	AirConditioning *bool `json:"airConditioning,omitempty"` // if AirConditioning exists
	Alarm           *bool `json:"alarm,omitempty"`           // if Alarm system exists
	Tv              *bool `json:"tv,omitempty"`              // if Tv exists
	Internet        *bool `json:"internet,omitempty"`        // if Internet exists
	DishWasher      *bool `json:"dishWasher,omitempty"`      // if DishWasher exists
	WashingMachine  *bool `json:"washingMachine,omitempty"`  // if WashingMachine exists
	Fridge          *bool `json:"fridge,omitempty"`          // if Fridge exists
	Stove           *bool `json:"stove,omitempty"`           // if Stove exists
}

// Interior details
type Interior struct {
	Structure         *AdRealtyStructureType         `json:"structure,omitempty"`      // internal Structure
	NumberOfBedrooms  *AdRealtyNumberOfBedroomsType  `json:"numOfBedrooms,omitempty"`  // NumberOfBedrooms number of bedrooms
	NumberOfBathrooms *AdRealtyNumberOfBathroomsType `json:"numOfBathrooms,omitempty"` // NumberOfBathrooms number of bathrooms
	Area              *Area                          `json:"area,omitempty"`           // Area details
	Equipment         *Equipment                     `json:"equipment,omitempty"`      // Equipment details
	RenovationType    *AdRealtyRenovationType        `json:"renovationType,omitempty"` // RenovationType type of renovation
	Furnished         *FurnishedType                 `json:"furnished,omitempty"`      // if Furnished
	Floor             *uint32                        `json:"floor,omitempty"`          // Floor
	Ceiling           *float32                       `json:"ceiling,omitempty"`        // Ceiling height in meters
	WinView           *AdRealtyWindowViewType        `json:"winView,omitempty"`        // WinView windows view
}

// Object details
type Object struct {
	Type     AdRealtyType `json:"type,omitempty"`     // Type of object
	Location *Location    `json:"location,omitempty"` // Location details
	Building *Building    `json:"building,omitempty"` // Building details
	Interior *Interior    `json:"interior,omitempty"` // Interior details
}

// SalePrice price of selling
type SalePrice struct {
	Cur    Currency `json:"cur,omitempty"`    // Cur selling currency
	Amount float64  `json:"amount,omitempty"` // selling Amount
}

// SaleTerms selling terms
type SaleTerms struct {
	Used  bool       `json:"used,omitempty"`  // if realty is Used
	Price *SalePrice `json:"price,omitempty"` // selling Price
}

// Realty details
type Realty struct {
	Rent   *RentTerms `json:"rent,omitempty"`   // populated for Rent ads
	Sale   *SaleTerms `json:"sale,omitempty"`   // populated for Sale ads
	Object *Object    `json:"object,omitempty"` // Object details
}

// AdImage is image details
type AdImage struct {
	FileId string `json:"fileId,omitempty"` // FileId is file link
	IsMain bool   `json:"isMain,omitempty"` // if image is IsMain
}

type AdDetails struct {
	Participants *Participants `json:"participants,omitempty"` // Participants details
	Description  *Description  `json:"description,omitempty"`  // Description details
	Realty       *Realty       `json:"realty,omitempty"`       // Realty details
	Images       []*AdImage    `json:"images,omitempty"`       // Images list
	Tags         []string      `json:"tags,omitempty"`         // Tags to search by
}

// Advertisement details
type Advertisement struct {
	Id          string     // Id unique identifier
	Code        string     // Code user friendly
	Status      string     // Status ads status
	SubStatus   string     // SubStatus ads sub status
	Type        string     // Type ads type
	SubType     string     // SubType ads subtype
	Details     *AdDetails // Details ads details
	ActivatedAt *time.Time // ActivatedAt when activated
	ClosedAt    *time.Time // ClosedAt when closed
}

// AdsSearchRequest search request
type AdsSearchRequest struct {
	kit.PagingRequest          // paging
	FullText          string   // FullText full text search
	Statuses          []string // Statuses search by statuses
	SubStatuses       []string // SubStatuses search by sub statuses
	Types             []string // Types search by types
	SubTypes          []string // SubTypes search by subtypes
}

// AdsInSearch represents ads in search response
type AdsInSearch struct {
	Id               string
	Code             string
	ShortDescription string
	Status           string
	SubStatus        string
	Type             string
	SubType          string
	Images           []*AdImage
}

type AdsSearchResponse struct {
	kit.PagingResponse                // paging
	Items              []*AdsInSearch // Items found items
}

// AdvertisementService provides ads functions
type AdvertisementService interface {
	// Create creates a new ads
	Create(ctx context.Context, ads *Advertisement) (*Advertisement, error)
	// Update updates an ads
	Update(ctx context.Context, ads *Advertisement) (*Advertisement, error)
	// Delete deletes an ads
	Delete(ctx context.Context, adsId string) error
	// SetStatus sets status
	SetStatus(ctx context.Context, adsId, status, subStatus string) (*Advertisement, error)
	// Get retrieves ads by id
	Get(ctx context.Context, adsId string) (*Advertisement, error)
	// Search searches ads by criteria
	Search(ctx context.Context, rq *AdsSearchRequest) (*AdsSearchResponse, error)
}

// AdvertisementStorage storage repository
type AdvertisementStorage interface {
	// Create creates a new ads
	Create(ctx context.Context, ads *Advertisement) error
	// Update updates an ads
	Update(ctx context.Context, ads *Advertisement) error
	// Delete deletes an ads
	Delete(ctx context.Context, adsId string) error
	// Get retrieves ads by id
	Get(ctx context.Context, adsId string) (*Advertisement, error)
	// Search searches ads by criteria
	Search(ctx context.Context, rq *AdsSearchRequest) (*AdsSearchResponse, error)
	// GetCode generates code
	GetCode(ctx context.Context) (string, error)
}
