package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/aidenappl/rootedparser/structs"
)

func main() {
	file, err := os.Open("./index_2025.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	finishedRecords := []structs.FullFilingRecord{}

	for _, row := range records[1:] {

		fmt.Println("Processing record:", row[5])

		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		xmlFileName := fmt.Sprintf("%s/%s/%s_public.xml", wd, row[9], row[8])
		xmlFile, err := os.Open(xmlFileName)
		if err != nil {
			fmt.Printf("Error opening XML file %s: %v\n", xmlFileName, err)
			continue
		}

		defer xmlFile.Close()

		var returnData structs.Return
		decoder := xml.NewDecoder(xmlFile)
		err = decoder.Decode(&returnData)
		if err != nil {
			fmt.Printf("Error decoding XML file %s: %v\n", xmlFileName, err)
			continue
		}

		fullRecord := structs.FullFilingRecord{
			EIN:        returnData.ReturnHeader.Filer.EIN,
			Name:       returnData.ReturnHeader.Filer.BusinessName.BusinessNameLine1Txt,
			DLN:        row[7],
			ObjectID:   row[8],
			XMLBatchID: row[9],
			Location: structs.Address{
				AddressLine1: returnData.ReturnHeader.Filer.USAddress.AddressLine1Txt,
				City:         returnData.ReturnHeader.Filer.USAddress.CityNm,
				State:        returnData.ReturnHeader.Filer.USAddress.StateAbbreviationCd,
				ZIPCode:      returnData.ReturnHeader.Filer.USAddress.ZIPCd,
			},
			Officers: []structs.Officer{},
		}

		if returnData.ReturnData.IRS990.FormationYear != 0 {
			fullRecord.Form990 = &structs.IRSForm990{
				PrincipalOfficerName: returnData.ReturnHeader.BusinessOfficerGrp[0].PersonNm,
				PrincipalOfficerAddress: structs.Address{
					AddressLine1: returnData.ReturnHeader.Filer.USAddress.AddressLine1Txt,
					City:         returnData.ReturnHeader.Filer.USAddress.CityNm,
					State:        returnData.ReturnHeader.Filer.USAddress.StateAbbreviationCd,
					ZIPCode:      returnData.ReturnHeader.Filer.USAddress.ZIPCd,
				},
				FormationYear:                   returnData.ReturnHeader.TaxYr,
				GrossReceiptsAmount:             returnData.ReturnData.IRS990.GrossReceipts,
				WebsiteAddress:                  returnData.ReturnData.IRS990.Website,
				MissionDescription:              returnData.ReturnData.IRS990.Mission,
				TotalAssetsEndOfYearAmount:      returnData.ReturnData.IRS990.TotalAssetsEOYAmt,
				TotalLiabilitiesEndOfYearAmount: returnData.ReturnData.IRS990.TotalLiabilitiesEOYAmt,
			}
		}

		if returnData.ReturnData.IRS990EZ.GrossReceiptsAmt != 0 {
			fullRecord.Form990EZ = &structs.IRS990EZ{
				GrossReceiptsAmt:          returnData.ReturnData.IRS990EZ.GrossReceiptsAmt,
				TotalRevenueAmt:           returnData.ReturnData.IRS990EZ.TotalRevenueAmt,
				TotalExpensesAmt:          returnData.ReturnData.IRS990EZ.TotalExpensesAmt,
				ExcessOrDeficitForYearAmt: returnData.ReturnData.IRS990EZ.ExcessOrDeficitForYearAmt,
				PrimaryExemptPurpose:      returnData.ReturnData.IRS990EZ.PrimaryExemptPurpose,
				President:                 returnData.ReturnData.IRS990EZ.President,

				PresidentHours: returnData.ReturnData.IRS990EZ.PresidentHours,
				Website:        returnData.ReturnData.IRS990EZ.Website,
				LocationCity:   returnData.ReturnData.IRS990EZ.LocationCity,
				LocationState:  returnData.ReturnData.IRS990EZ.LocationState,
			}
		}

		// Add officers from the XML data
		for _, officer := range returnData.ReturnHeader.BusinessOfficerGrp {
			fullRecord.Officers = append(fullRecord.Officers, structs.Officer{
				PersonName:                 officer.PersonNm,
				PersonTitle:                officer.PersonTitleTxt,
				PhoneNumber:                officer.PhoneNum,
				SignatureDate:              officer.SignatureDt,
				DiscussWithPaidPreparerInd: officer.DiscussWithPaidPreparerInd,
			})
		}

		finishedRecords = append(finishedRecords, fullRecord)

		xmlFile.Close()
		fmt.Println("Finished processing record:", row[0])
		fmt.Println("--------------------------------------------------")
	}

	outputFile, err := os.Create("finished_records.json")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	encoder := json.NewEncoder(outputFile)
	err = encoder.Encode(finishedRecords)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finished records saved to finished_records.json")
	fmt.Println("All records processed successfully.")
	fmt.Println("Total records processed:", len(finishedRecords))
	fmt.Println("Exiting program.")
	fmt.Println("Goodbye!")
	fmt.Println("--------------------------------------------------")
}
