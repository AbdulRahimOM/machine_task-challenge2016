package regions

import (
	"challenge16/utils"
)

const (
	filePath = "cities.csv"
)

type (
	countryData struct {
		Name      string
		Provinces map[string]provinceData
	}

	provinceData struct {
		Name   string
		Cities map[string]string
	}
)

var Countries = make(map[string]countryData)

func LoadDataIntoMap(csvFilePath string) error {
	// Load data from CSV into the countries map

	datas, err := utils.ParseCSV(csvFilePath)
	if err != nil {
		return err
	}

	for _, data := range datas {
		// Add data to the map
		if _, ok := Countries[data.CountryCode]; !ok {
			Countries[data.CountryCode] = countryData{
				Name: data.CountryName,
				Provinces: map[string]provinceData{
					data.ProvinceCode: {
						Name: data.ProvinceName,
						Cities: map[string]string{
							data.CityCode: data.CityName,
						},
					},
				},
			}
			continue
		}
		if _, ok := Countries[data.CountryCode].Provinces[data.ProvinceCode]; !ok {
			Countries[data.CountryCode].Provinces[data.ProvinceCode] = provinceData{
				Name:   data.ProvinceName,
				Cities: map[string]string{data.CityCode: data.CityName},
			}
			continue
		}
		Countries[data.CountryCode].Provinces[data.ProvinceCode].Cities[data.CityCode] = data.CityName
	}

	return nil
}
