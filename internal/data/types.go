package data

type permissionData struct {
	/*
		heirarchy...
		country->(if country is mentioned     ): -(excludedProvinces) -(excludedCities)
		country->(if country is not mentioned ): +(includedProvinces - excludedCities) + (includedCities)
	*/
	includedCountries map[string]bool
	includedProvinces map[string]map[string]bool
	excludedProvinces map[string]map[string]bool
	includedCities    map[string]map[string]map[string]bool
	excludedCities    map[string]map[string]map[string]bool
}

func (src permissionData) copyPermissionData() permissionData {
	dst := permissionData{
		includedCountries: make(map[string]bool),
		includedProvinces: make(map[string]map[string]bool),
		excludedProvinces: make(map[string]map[string]bool),
		includedCities:    make(map[string]map[string]map[string]bool),
		excludedCities:    make(map[string]map[string]map[string]bool),
	}

	// Copy includedCountries
	for k, v := range src.includedCountries {
		dst.includedCountries[k] = v
	}

	// Copy includedProvinces
	for k, v := range src.includedProvinces {
		dst.includedProvinces[k] = make(map[string]bool)
		for k2, v2 := range v {
			dst.includedProvinces[k][k2] = v2
		}
	}

	// Copy excludedProvinces
	for k, v := range src.excludedProvinces {
		dst.excludedProvinces[k] = make(map[string]bool)
		for k2, v2 := range v {
			dst.excludedProvinces[k][k2] = v2
		}
	}

	// Copy includedCities
	for country, provinces := range src.includedCities {
		dst.includedCities[country] = make(map[string]map[string]bool)
		for province, cities := range provinces {
			dst.includedCities[country][province] = make(map[string]bool)
			for city, v := range cities {
				dst.includedCities[country][province][city] = v
			}
		}
	}

	// Copy excludedCities
	for country, provinces := range src.excludedCities {
		dst.excludedCities[country] = make(map[string]map[string]bool)
		for province, cities := range provinces {
			dst.excludedCities[country][province] = make(map[string]bool)
			for city, v := range cities {
				dst.excludedCities[country][province][city] = v
			}
		}
	}

	return dst
}
