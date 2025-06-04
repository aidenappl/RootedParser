package structs

type FullFilingRecord struct {
	EIN        string      `json:"ein"`
	Name       string      `json:"name"`
	DLN        string      `json:"dln"`
	ObjectID   string      `json:"object_id"`
	XMLBatchID string      `json:"xml_batch_id"`
	Location   Address     `json:"location"`
	People     []People    `json:"people"`
	Form990    *IRSForm990 `json:"form_990"`
	Form990EZ  *IRS990EZ   `json:"form_990_ez"`
}

type People struct {
	PersonName   string   `json:"person_name"`
	PersonTitle  string   `json:"person_title"`
	PhoneNumber  *string  `json:"phone_number"`
	AverageHours *float64 `json:"average_hours"`
	Bookkeeper   bool     `json:"bookkeeper"`
	Compensation *int     `json:"compensation"`
	Address      *Address `json:"address"`
}

type IRS990EZ struct {
	GrossReceiptsAmt          int    `json:"gross_receipts_amt"`
	TotalRevenueAmt           int    `json:"total_revenue_amt"`
	TotalExpensesAmt          int    `json:"total_expenses_amt"`
	ExcessOrDeficitForYearAmt int    `json:"excess_or_deficit_for_year_amt"`
	PrimaryExemptPurpose      string `json:"primary_exempt_purpose_txt"`
	Website                   string `json:"website"`
}

type IRSForm990 struct {
	PrincipalOfficerName            string  `json:"principal_officer_name"`
	PrincipalOfficerAddress         Address `json:"principal_officer_address"`
	GrossReceiptsAmount             int     `json:"gross_receipts_amount"`
	WebsiteAddress                  string  `json:"website_address"`
	MissionDescription              string  `json:"mission_description"`
	FormationYear                   int     `json:"formation_year"`
	TotalAssetsEndOfYearAmount      int     `json:"total_assets_end_of_year_amount"`
	TotalLiabilitiesEndOfYearAmount int     `json:"total_liabilities_end_of_year_amount"`
}

type Address struct {
	AddressLine1 string `json:"address_line_1"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZIPCode      string `json:"zip_code"`
}
