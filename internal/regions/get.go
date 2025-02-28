package regions

import (
	"challenge16/internal/response"
	"fmt"
)

type regionInfo struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func GetCountries() response.Response {
	countries := make([]regionInfo, 0, len(Countries))
	for code, country := range Countries {
		countries = append(countries, regionInfo{
			Name: country.Name,
			Code: code,
		})
	}
	return response.CreateSuccess(200, "SUCCESS", map[string]interface{}{
		"countries": countries,
	})
}

func GetProvincesInCountry(countryCode string) response.Response {
	exists := checkCountry(countryCode)
	if !exists {
		return response.CreateError(404, "COUNTRY_NOT_FOUND", fmt.Errorf("Country not found: %s", countryCode))
	}

	country := Countries[countryCode]
	provinces := make([]regionInfo, 0, len(country.Provinces))
	for code, province := range country.Provinces {
		provinces = append(provinces, regionInfo{
			Name: province.Name,
			Code: code,
		})
	}
	return response.CreateSuccess(200, "SUCCESS", map[string]interface{}{
		"provinces": provinces,
	})
}

func GetCitiesInProvince(countryCode, provinceCode string) response.Response {
	exists := checkProvince(countryCode, provinceCode)
	if !exists {
		return response.CreateError(404, "PROVINCE_NOT_FOUND", fmt.Errorf("Province not found: %s-%s", countryCode, provinceCode))
	}

	province := Countries[countryCode].Provinces[provinceCode]
	cities := make([]regionInfo, 0, len(province.Cities))
	for code, name := range province.Cities {
		cities = append(cities, regionInfo{
			Name: name,
			Code: code,
		})
	}
	return response.CreateSuccess(200, "SUCCESS", map[string]interface{}{
		"cities": cities,
	})
}
