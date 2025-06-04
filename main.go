package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/aidenappl/rootedparser/structs"
	"github.com/schollz/progressbar/v3"
)

func main() {
	fmt.Println("Starting to process records...")

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
	errorLog := make([]string, 0)
	fmt.Println("Total records to process:", len(records)-1)

	// wait for user input to continue
	fmt.Println("Press Enter to continue...")
	var input string
	fmt.Scanln(&input)

	// check that user pressed enter
	if input != "" {
		fmt.Println("Exiting program.")
		return
	}

	fmt.Println("Processing records...")
	fmt.Println()
	bar := progressbar.Default(int64(len(records) - 1))

	for _, row := range records[1:] {

		bar.Add(1)

		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		xmlFileName := fmt.Sprintf("%s/%s/%s_public.xml", wd, row[9], row[8])
		xmlFile, err := os.Open(xmlFileName)
		if err != nil {
			fmt.Printf("Error opening XML file %s: %v\n", xmlFileName, err)
			errorLog = append(errorLog, fmt.Sprintf("Error opening XML file %s: %v", xmlFileName, err))
			continue
		}

		defer xmlFile.Close()

		var returnData structs.Return
		decoder := xml.NewDecoder(xmlFile)
		err = decoder.Decode(&returnData)
		if err != nil {
			fmt.Printf("Error decoding XML file %s: %v\n", xmlFileName, err)
			errorLog = append(errorLog, fmt.Sprintf("Error decoding XML file %s: %v", xmlFileName, err))
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
			People: []structs.People{},
		}
		// -
		// Checking if the returnData has the IRS990 Attached
		// Handling 990 data
		// -
		if returnData.ReturnData.IRS990 != nil {
			if returnData.ReturnData.IRS990.People != nil {
				for _, person := range returnData.ReturnData.IRS990.People {
					fullRecord.People = append(fullRecord.People, structs.People{
						PersonName:   person.Name,
						PersonTitle:  person.Title,
						AverageHours: &person.AverageHoursPerWeek,
						Compensation: &person.CompFromOrg,
					})
				}
			}

			// Adding the principal officer from the return header to the people slice
			fullRecord.People = append(fullRecord.People, structs.People{
				PersonName:  returnData.ReturnHeader.BusinessOfficerGrp[0].PersonNm,
				PersonTitle: returnData.ReturnHeader.BusinessOfficerGrp[0].PersonTitleTxt,
				PhoneNumber: &returnData.ReturnHeader.BusinessOfficerGrp[0].PhoneNum,
				Address: &structs.Address{
					AddressLine1: returnData.ReturnHeader.Filer.USAddress.AddressLine1Txt,
					City:         returnData.ReturnHeader.Filer.USAddress.CityNm,
					State:        returnData.ReturnHeader.Filer.USAddress.StateAbbreviationCd,
					ZIPCode:      returnData.ReturnHeader.Filer.USAddress.ZIPCd,
				},
			})

			// Adding the IRS990 data to the fullRecord
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

		// -
		// Checking if the returnData has the IRS990EZ Attached
		// Handling 990EZ data
		// -

		if returnData.ReturnData.IRS990EZ != nil {
			fullRecord.Form990EZ = &structs.IRS990EZ{
				GrossReceiptsAmt:          returnData.ReturnData.IRS990EZ.GrossReceiptsAmt,
				TotalRevenueAmt:           returnData.ReturnData.IRS990EZ.TotalRevenueAmt,
				TotalExpensesAmt:          returnData.ReturnData.IRS990EZ.TotalExpensesAmt,
				ExcessOrDeficitForYearAmt: returnData.ReturnData.IRS990EZ.ExcessOrDeficitForYearAmt,
				PrimaryExemptPurpose:      returnData.ReturnData.IRS990EZ.PrimaryExemptPurpose,
				Website:                   returnData.ReturnData.IRS990EZ.Website,
			}

			// Add bookkeeper information to the people slice
			if returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.PersonName != "" {
				//  Check if person is already added to avoid duplicates
				alreadyExists := false
				for idx, person := range fullRecord.People {
					if person.PersonName == returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.PersonName {
						alreadyExists = true
						// Update existing person's details
						fullRecord.People[idx].Address = &structs.Address{
							AddressLine1: returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.Address.AddressLine1Txt,
							City:         returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.Address.CityNm,
							State:        returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.Address.StateAbbreviationCd,
							ZIPCode:      returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.Address.ZIPCd,
						}
						fullRecord.People[idx].Bookkeeper = true
						fullRecord.People[idx].PhoneNumber = &returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.PhoneNumber
						break
					}
				}
				if !alreadyExists {
					if returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.Address != nil && returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.Address.AddressLine1Txt != "" {
						fullRecord.People = append(fullRecord.People, structs.People{
							PersonName: returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.PersonName,
							Address: &structs.Address{
								AddressLine1: returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.Address.AddressLine1Txt,
								City:         returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.Address.CityNm,
								State:        returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.Address.StateAbbreviationCd,
								ZIPCode:      returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.Address.ZIPCd,
							},
							Bookkeeper:  true,
							PhoneNumber: &returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.PhoneNumber,
						})
					} else {
						fullRecord.People = append(fullRecord.People, structs.People{
							PersonName:  returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.PersonName,
							Bookkeeper:  true,
							PhoneNumber: &returnData.ReturnData.IRS990EZ.BooksInCareOfDetail.PhoneNumber,
						})
					}

				}
			}

			// Adding the officers from the IRS990EZ to the people slice
			if len(returnData.ReturnData.IRS990EZ.OfficerDirectorTrusteeEmplGrp) != 0 {
				for _, officer := range returnData.ReturnData.IRS990EZ.OfficerDirectorTrusteeEmplGrp {
					// check if person is already added to avoid duplicates
					alreadyExists := false
					for idx, person := range fullRecord.People {
						if person.PersonName == officer.Name {
							// If the person already exists, update their compensation and average hours
							fullRecord.People[idx].Compensation = &officer.Compensation
							fullRecord.People[idx].AverageHours = &officer.AverageHoursPerWk
							fullRecord.People[idx].PersonTitle = officer.Title
							if officer.Address != nil {
								fullRecord.People[idx].Address = &structs.Address{
									AddressLine1: officer.Address.AddressLine1Txt,
									City:         officer.Address.CityNm,
									State:        officer.Address.StateAbbreviationCd,
									ZIPCode:      officer.Address.ZIPCd,
								}
							}
							alreadyExists = true
							break
						}
					}
					if !alreadyExists {
						fullRecord.People = append(fullRecord.People, structs.People{
							PersonName:   officer.Name,
							PersonTitle:  officer.Title,
							AverageHours: &officer.AverageHoursPerWk,
							Compensation: &officer.Compensation,
						})
					}
				}
			}
		}

		// -
		// Checking if the returnData has the IRS990PF Attached
		// Handling 990PF data
		// -
		if returnData.ReturnData.IRS990PF != nil && len(returnData.ReturnData.IRS990PF.Officers) > 0 {
			for _, pfOfficer := range returnData.ReturnData.IRS990PF.Officers {

				// check if person is already added to avoid duplicates
				alreadyExists := false
				for idx, person := range fullRecord.People {
					if person.PersonName == pfOfficer.Name {
						alreadyExists = true
						// Update existing person's details
						fullRecord.People[idx].AverageHours = &pfOfficer.AverageHoursPerWk
						fullRecord.People[idx].Compensation = &pfOfficer.Compensation
						if pfOfficer.Address != nil {
							fullRecord.People[idx].Address = &structs.Address{
								AddressLine1: pfOfficer.Address.AddressLine1Txt,
								City:         pfOfficer.Address.CityNm,
								State:        pfOfficer.Address.StateAbbreviationCd,
								ZIPCode:      pfOfficer.Address.ZIPCd,
							}
						}
						break
					}
				}
				if alreadyExists {
					continue
				}

				// check if the officer has a valid address
				if pfOfficer.Address != nil && pfOfficer.Address.AddressLine1Txt != "" {
					fullRecord.People = append(fullRecord.People, structs.People{
						PersonName:   pfOfficer.Name,
						PersonTitle:  pfOfficer.Title,
						AverageHours: &pfOfficer.AverageHoursPerWk,
						Compensation: &pfOfficer.Compensation,
						Address: &structs.Address{
							AddressLine1: pfOfficer.Address.AddressLine1Txt,
							City:         pfOfficer.Address.CityNm,
							State:        pfOfficer.Address.StateAbbreviationCd,
							ZIPCode:      pfOfficer.Address.ZIPCd,
						},
					})
				} else {
					fullRecord.People = append(fullRecord.People, structs.People{
						PersonName:   pfOfficer.Name,
						PersonTitle:  pfOfficer.Title,
						AverageHours: &pfOfficer.AverageHoursPerWk,
						Compensation: &pfOfficer.Compensation,
					})
				}
			}
		}

		// -
		// Checking if the returnData has the Business Officer Groups Attached
		// Handling Business Officer Group data
		// -
		for _, officer := range returnData.ReturnHeader.BusinessOfficerGrp {
			//  check if person is already added to avoid duplicates
			alreadyExists := false
			for idx, person := range fullRecord.People {
				if person.PersonName == officer.PersonNm {
					// Update existing person's details
					fullRecord.People[idx].PersonTitle = officer.PersonTitleTxt
					fullRecord.People[idx].PhoneNumber = &officer.PhoneNum
					alreadyExists = true
					break
				}
			}
			if alreadyExists {
				continue
			}
			fullRecord.People = append(fullRecord.People, structs.People{
				PersonName:  officer.PersonNm,
				PersonTitle: officer.PersonTitleTxt,
				PhoneNumber: &officer.PhoneNum,
			})
		}

		// Validate people records, check for duplicates and empty names
		peopleMap := make(map[string]structs.People)
		for _, person := range fullRecord.People {
			// Skip if person name is empty
			if person.PersonName == "" {
				continue
			}
			// Check if the person already exists in the map
			if existingPerson, exists := peopleMap[person.PersonName]; exists {
				// If the person already exists, update their details
				if person.PersonTitle != "" {
					existingPerson.PersonTitle = person.PersonTitle
				}
				if person.PhoneNumber != nil {
					existingPerson.PhoneNumber = person.PhoneNumber
				}
				if person.AverageHours != nil {
					existingPerson.AverageHours = person.AverageHours
				}
				if person.Compensation != nil {
					existingPerson.Compensation = person.Compensation
				}
				if person.Address != nil {
					existingPerson.Address = person.Address
				}
				// Update the map with the existing person
				peopleMap[person.PersonName] = existingPerson
			} else {
				// If the person does not exist, add them to the map
				peopleMap[person.PersonName] = person
			}
		}
		// Convert the map back to a slice
		fullRecord.People = make([]structs.People, 0, len(peopleMap))
		for _, person := range peopleMap {
			fullRecord.People = append(fullRecord.People, person)
		}

		// Remove any empty addresses
		for i := len(fullRecord.People) - 1; i >= 0; i-- {
			if fullRecord.People[i].Address != nil && fullRecord.People[i].Address.AddressLine1 == "" &&
				fullRecord.People[i].Address.City == "" &&
				fullRecord.People[i].Address.State == "" && fullRecord.People[i].Address.ZIPCode == "" {
				fullRecord.People[i].Address = nil
			}
		}
		// Remove any empty phone numbers
		for i := len(fullRecord.People) - 1; i >= 0; i-- {
			if fullRecord.People[i].PhoneNumber != nil && *fullRecord.People[i].PhoneNumber == "" {
				fullRecord.People[i].PhoneNumber = nil
			}
		}

		// Send the full record to the finished records slice
		finishedRecords = append(finishedRecords, fullRecord)

		xmlFile.Close()
	}

	// Log errors if any
	if len(errorLog) > 0 {
		fmt.Println("Errors encountered during processing:")
		outputErrFile, err := os.Create("error_log.txt")
		if err != nil {
			panic(err)
		}
		defer outputErrFile.Close()
		for _, errMsg := range errorLog {
			fmt.Println(errMsg)
			_, err := outputErrFile.WriteString(errMsg + "\n")
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("Error log saved to error_log.txt")
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
	fmt.Println()
	fmt.Println("Finished records saved to finished_records.json")
	fmt.Println("All records processed successfully.")
	fmt.Println()
	fmt.Println("Records requested:", len(records)-1)
	fmt.Println("Total records processed:", len(finishedRecords))
	fmt.Println("Total errors encountered:", len(errorLog))
	fmt.Println("Success rate: ", float64(len(finishedRecords))/float64(len(records)-1)*100, "%")
	fmt.Println()
	fmt.Println("Exiting program.")
	fmt.Println("Goodbye!")
	fmt.Println("--------------------------------------------------")
}
