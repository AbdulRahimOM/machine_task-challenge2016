package utils

import (
	"encoding/csv"
	"os"
)

type Data struct {
	CityCode     string
	CityName     string
	ProvinceCode string
	ProvinceName string
	CountryCode  string
	CountryName  string
}

func ParseCSV(filename string) ([]Data, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var dataList []Data
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		dataList = append(dataList, Data{
			CityCode:     record[0],
			ProvinceCode: record[1],
			CountryCode:  record[2],
			CityName:     record[3],
			ProvinceName: record[4],
			CountryName:  record[5],
		})
	}

	return dataList, nil
}
