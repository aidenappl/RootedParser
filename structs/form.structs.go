package structs

type FullFilingRecord struct {
	EIN        string      `json:"ein"`
	Name       string      `json:"name"`
	DLN        string      `json:"dln"`
	ObjectID   string      `json:"object_id"`
	XMLBatchID string      `json:"xml_batch_id"`
	Location   Address     `json:"location"`
	Officers   []Officer   `json:"officers"`
	Form990    *IRSForm990 `json:"form_990"`
	Form990EZ  *IRS990EZ   `json:"form_990_ez"`
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

type Officer struct {
	PersonName                 string `json:"person_name"`
	PersonTitle                string `json:"person_title"`
	PhoneNumber                string `json:"phone_number"`
	SignatureDate              string `json:"signature_date"`
	DiscussWithPaidPreparerInd string `json:"discuss_with_paid_preparer_ind"`
}

type Return struct {
	ReturnHeader ReturnHeader `xml:"ReturnHeader"`
	ReturnData   ReturnData   `xml:"ReturnData"`
}

type ReturnHeader struct {
	ReturnTs           string            `xml:"ReturnTs"`
	TaxPeriodEndDt     string            `xml:"TaxPeriodEndDt"`
	TaxPeriodBeginDt   string            `xml:"TaxPeriodBeginDt"`
	ReturnTypeCd       string            `xml:"ReturnTypeCd"`
	TaxYr              int               `xml:"TaxYr"`
	Filer              Filer             `xml:"Filer"`
	BusinessOfficerGrp []BusinessOfficer `xml:"BusinessOfficerGrp"`
	PreparerPersonGrp  PreparerPerson    `xml:"PreparerPersonGrp"`
}

type BusinessName struct {
	BusinessNameLine1Txt string `xml:"BusinessNameLine1Txt"`
}

type USAddress struct {
	AddressLine1Txt     string `xml:"AddressLine1Txt"`
	CityNm              string `xml:"CityNm"`
	StateAbbreviationCd string `xml:"StateAbbreviationCd"`
	ZIPCd               string `xml:"ZIPCd"`
}

type Filer struct {
	EIN          string       `xml:"EIN"`
	BusinessName BusinessName `xml:"BusinessName"`
	PhoneNum     string       `xml:"PhoneNum"`
	USAddress    USAddress    `xml:"USAddress"`
}

type BusinessOfficer struct {
	PersonNm                   string `xml:"PersonNm"`
	PersonTitleTxt             string `xml:"PersonTitleTxt"`
	PhoneNum                   string `xml:"PhoneNum"`
	SignatureDt                string `xml:"SignatureDt"`
	DiscussWithPaidPreparerInd string `xml:"DiscussWithPaidPreparerInd"`
}

type PreparerPerson struct {
	PreparerPersonNm string `xml:"PreparerPersonNm"`
	PTIN             string `xml:"PTIN"`
	PhoneNum         string `xml:"PhoneNum"`
	PreparationDt    string `xml:"PreparationDt"`
}

type ReturnData struct {
	DocumentCnt int        `xml:"documentCnt,attr"`
	IRS990EZ    IRS990EZ   `xml:"IRS990EZ"`
	IRS990      IRS990Form `xml:"IRS990"`
}

type IRS990EZ struct {
	GrossReceiptsAmt          int     `xml:"GrossReceiptsAmt"`
	TotalRevenueAmt           int     `xml:"TotalRevenueAmt"`
	TotalExpensesAmt          int     `xml:"TotalExpensesAmt"`
	ExcessOrDeficitForYearAmt int     `xml:"ExcessOrDeficitForYearAmt"`
	PrimaryExemptPurpose      string  `xml:"PrimaryExemptPurposeTxt"`
	President                 string  `xml:"OfficerDirectorTrusteeEmplGrp>PersonNm"`
	PresidentHours            float64 `xml:"OfficerDirectorTrusteeEmplGrp>AverageHrsPerWkDevotedToPosRt"`
	Website                   string  `xml:"WebsiteAddressTxt"`
	LocationCity              string  `xml:"BooksInCareOfDetail>USAddress>CityNm"`
	LocationState             string  `xml:"BooksInCareOfDetail>USAddress>StateAbbreviationCd"`
}

type IRS990Form struct {
	DoingBusinessAs        string `xml:"DoingBusinessAsName>BusinessNameLine1Txt"`
	Website                string `xml:"WebsiteAddressTxt"`
	GrossReceipts          int    `xml:"GrossReceiptsAmt"`
	FormationYear          int    `xml:"FormationYr"`
	State                  string `xml:"LegalDomicileStateCd"`
	Mission                string `xml:"ActivityOrMissionDesc"`
	TotalVolunteers        int    `xml:"TotalVolunteersCnt"`
	TotalRevenue           int    `xml:"CYTotalRevenueAmt"`
	TotalExpenses          int    `xml:"CYTotalExpensesAmt"`
	RevenueLessExpenses    int    `xml:"CYRevenuesLessExpensesAmt"`
	Contributions          int    `xml:"CYContributionsGrantsAmt"`
	MembershipDues         int    `xml:"MembershipDuesAmt"`
	ProgramServiceRevenue  int    `xml:"CYProgramServiceRevenueAmt"`
	InvestmentIncome       int    `xml:"CYInvestmentIncomeAmt"`
	OtherRevenue           int    `xml:"CYOtherRevenueAmt"`
	TotalAssetsEOYAmt      int    `xml:"TotalAssetsEOYAmt"`
	TotalLiabilitiesEOYAmt int    `xml:"TotalLiabilitiesEOYAmt"`
}
