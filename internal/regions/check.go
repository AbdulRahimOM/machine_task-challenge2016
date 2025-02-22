package regions

func CheckCountry(countryCode string) bool {
	_, ok := Countries[countryCode]
	return ok
}

func CheckProvince(countryCode, provinceCode string) bool {
	if CheckCountry(countryCode) == false {
		return false
	}
	_, ok := Countries[countryCode].Provinces[provinceCode]
	return ok
}

func CheckCity(countryCode, provinceCode, cityCode string) bool {
	if CheckProvince(countryCode, provinceCode) == false {
		return false
	}
	_, ok := Countries[countryCode].Provinces[provinceCode].Cities[cityCode]
	return ok
}
