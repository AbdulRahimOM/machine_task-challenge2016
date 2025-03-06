package regions

type regionInfo struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func GetCountries() []regionInfo {
	countries := make([]regionInfo, 0, len(Countries))
	for code, country := range Countries {
		countries = append(countries, regionInfo{
			Name: country.Name,
			Code: code,
		})
	}
	return countries
}

func GetProvincesInCountry(countryCode string) []regionInfo {
	if !CheckCountry(countryCode) {
		return nil
	}
	country := Countries[countryCode]
	if country.Provinces == nil {
		return nil
	}
	provinces := make([]regionInfo, 0, len(country.Provinces))
	for code, province := range country.Provinces {
		provinces = append(provinces, regionInfo{
			Name: province.Name,
			Code: code,
		})
	}
	return provinces
}

func GetCitiesInProvince(countryCode, provinceCode string) []regionInfo {
	if !CheckProvince(countryCode, provinceCode) {
		return nil
	}
	province := Countries[countryCode].Provinces[provinceCode]
	cities := make([]regionInfo, 0, len(province.Cities))
	for code, name := range province.Cities {
		cities = append(cities, regionInfo{
			Name: name,
			Code: code,
		})
	}
	return cities
}
