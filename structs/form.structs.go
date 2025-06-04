package structs

type Return struct {
	ReturnHeader ReturnHeader `xml:"ReturnHeader"`
	ReturnData   ReturnData   `xml:"ReturnData"`
}

type ReturnData struct {
	DocumentCnt int          `xml:"documentCnt,attr"`
	IRS990EZ    *IRS990EZ_FL `xml:"IRS990EZ"`
	IRS990      *IRS990Form  `xml:"IRS990"`
	IRS990PF    *IRS990PF_FL `xml:"IRS990PF"`
}

type IRS990EZ_FL struct {
	GrossReceiptsAmt              int                  `xml:"GrossReceiptsAmt"`
	TotalRevenueAmt               int                  `xml:"TotalRevenueAmt"`
	TotalExpensesAmt              int                  `xml:"TotalExpensesAmt"`
	ExcessOrDeficitForYearAmt     int                  `xml:"ExcessOrDeficitForYearAmt"`
	OfficerDirectorTrusteeEmplGrp []PFOfficer_FL       `xml:"OfficerDirectorTrusteeEmplGrp"`
	BooksInCareOfDetail           IRS990EZ_BooksInCare `xml:"BooksInCareOfDetail"`
	PrimaryExemptPurpose          string               `xml:"PrimaryExemptPurposeTxt"`
	Website                       string               `xml:"WebsiteAddressTxt"`
}

type IRS990EZ_BooksInCare struct {
	PersonName  string     `xml:"PersonNm"`
	Address     *USAddress `xml:"USAddress"`
	PhoneNumber string     `xml:"PhoneNum"`
}

type IRS990PF_FL struct {
	Officers []PFOfficer_FL `xml:"OfficerDirTrstKeyEmplInfoGrp>OfficerDirTrstKeyEmplGrp"`
}

type IRS990Form struct {
	DoingBusinessAs        string   `xml:"DoingBusinessAsName>BusinessNameLine1Txt"`
	Website                string   `xml:"WebsiteAddressTxt"`
	GrossReceipts          int      `xml:"GrossReceiptsAmt"`
	FormationYear          int      `xml:"FormationYr"`
	State                  string   `xml:"LegalDomicileStateCd"`
	Mission                string   `xml:"ActivityOrMissionDesc"`
	TotalVolunteers        int      `xml:"TotalVolunteersCnt"`
	TotalRevenue           int      `xml:"CYTotalRevenueAmt"`
	TotalExpenses          int      `xml:"CYTotalExpensesAmt"`
	RevenueLessExpenses    int      `xml:"CYRevenuesLessExpensesAmt"`
	Contributions          int      `xml:"CYContributionsGrantsAmt"`
	MembershipDues         int      `xml:"MembershipDuesAmt"`
	ProgramServiceRevenue  int      `xml:"CYProgramServiceRevenueAmt"`
	InvestmentIncome       int      `xml:"CYInvestmentIncomeAmt"`
	OtherRevenue           int      `xml:"CYOtherRevenueAmt"`
	TotalAssetsEOYAmt      int      `xml:"TotalAssetsEOYAmt"`
	TotalLiabilitiesEOYAmt int      `xml:"TotalLiabilitiesEOYAmt"`
	People                 []Person `xml:"Form990PartVIISectionAGrp"`
}

type PFOfficer_FL struct {
	Name              string     `xml:"PersonNm"`
	Title             string     `xml:"TitleTxt"`
	AverageHoursPerWk float64    `xml:"AverageHrsPerWkDevotedToPosRt"`
	Compensation      int        `xml:"CompensationAmt"`
	Address           *USAddress `xml:"USAddress,omitempty"`
}

type Person struct {
	Name                        string  `xml:"PersonNm"`
	Title                       string  `xml:"TitleTxt"`
	AverageHoursPerWeek         float64 `xml:"AverageHoursPerWeekRt"`
	AverageHoursRelatedOrgs     float64 `xml:"AverageHoursPerWeekRltdOrgRt"`
	IndividualTrusteeOrDirector string  `xml:"IndividualTrusteeOrDirectorInd"`
	Officer                     string  `xml:"OfficerInd"`
	CompFromOrg                 int     `xml:"ReportableCompFromOrgAmt"`
	CompFromRelatedOrgs         int     `xml:"ReportableCompFromRltdOrgAmt"`
	OtherCompensation           int     `xml:"OtherCompensationAmt"`
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

type Filer struct {
	EIN          string       `xml:"EIN"`
	BusinessName BusinessName `xml:"BusinessName"`
	PhoneNum     string       `xml:"PhoneNum"`
	USAddress    USAddress    `xml:"USAddress"`
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
