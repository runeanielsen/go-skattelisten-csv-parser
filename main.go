package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type company struct {
	Cvr           string
	Name          string
	Se            string
	IncomeYear    int
	CompanyType   string
	TaxableIncome string
	Deficit       int64
	CorporateTax  int64
}

func main() {
	source := flag.String("s", "", "Source csv file")
	destination := flag.String("d", "", "Destination JSON file.")
	flag.Parse()

	f, err := os.Open(*source)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	companies := []company{}
	for index, line := range text {
		// Skip first value since it is shit.
		if index == 0 {
			continue
		}

		splittedString := strings.Split(line, ",")
		company := createCompany(splittedString)
		companies = append(companies, company)
	}

	jsonCompanies, _ := json.Marshal(companies)
	ioutil.WriteFile(*destination, jsonCompanies, 0644)
}

func createCompany(attributes []string) company {
	incomeYearStr := strings.TrimSpace(attributes[3])
	if incomeYearStr == "" {
		incomeYearStr = "0"
	}
	incomeYear, err := strconv.Atoi(incomeYearStr)
	if err != nil {
		log.Fatalf("Could not parse income year %s", incomeYearStr)
	}

	deficitStr := strings.TrimSpace(attributes[9])
	if deficitStr == "" {
		deficitStr = "0"
	}
	deficit, err := strconv.ParseInt(deficitStr, 10, 64)
	if err != nil {
		log.Fatalf("Could not parse deficit '%s'", deficitStr)
	}

	corporateTaxStr := strings.TrimSpace(attributes[10])
	if corporateTaxStr == "" {
		corporateTaxStr = "0"
	}
	corporateTax, err := strconv.ParseInt(corporateTaxStr, 10, 64)
	if err != nil {
		log.Fatalf("Could not parse corporate tax '%s'", corporateTaxStr)
	}

	return company{
		Cvr:           attributes[0],
		Name:          attributes[1],
		Se:            attributes[2],
		IncomeYear:    incomeYear,
		CompanyType:   attributes[5],
		TaxableIncome: attributes[8],
		Deficit:       deficit,
		CorporateTax:  corporateTax,
	}
}
