package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
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
	companies := []company{}

	// Skip first line since it should not be parsed
	scanner.Scan()
	for scanner.Scan() {
		splittedString := strings.Split(scanner.Text(), ",")
		company, err := createCompany(splittedString)
		if err != nil {
			log.Fatal(err)
		}

		companies = append(companies, company)
	}

	jsonCompanies, err := json.Marshal(companies)
	if err != nil {
		log.Fatal("Could not parse companies to json")
	}
	ioutil.WriteFile(*destination, jsonCompanies, 0644)
}

func createCompany(attributes []string) (company, error) {
	incomeYearStr := strings.TrimSpace(attributes[3])
	if incomeYearStr == "" {
		incomeYearStr = "0"
	}
	incomeYear, err := strconv.Atoi(incomeYearStr)
	if err != nil {
		return company{}, fmt.Errorf("Could not parse incomeYear '%s'", incomeYearStr)
	}

	deficitStr := strings.TrimSpace(attributes[9])
	if deficitStr == "" {
		deficitStr = "0"
	}
	deficit, err := strconv.ParseInt(deficitStr, 10, 64)
	if err != nil {
		return company{}, fmt.Errorf("Could not parse deficit '%s'", deficitStr)
	}

	corporateTaxStr := strings.TrimSpace(attributes[10])
	if corporateTaxStr == "" {
		corporateTaxStr = "0"
	}
	corporateTax, err := strconv.ParseInt(corporateTaxStr, 10, 64)
	if err != nil {
		return company{}, fmt.Errorf("Could not parse corporate tax '%s'", corporateTaxStr)
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
	}, nil
}
