package dto

import "challenge16/internal/regions"

type (
	// Contract struct {
	// 	ParentDistributor string
	// 	SubDistributor    *string
	// 	IncludedRegions   []string
	// 	ExcludedRegions   []string
	// }

	Contract struct {
		/*
			heirarchy...
			country->(if country is mentioned     ): (-excludedProvinces) - (excludedCities)
			country->(if country is not mentioned ): ( includedProvinces - excludedCities) + (includedCities)
		*/

		ParentDistributor *string
		ContractRecipient string
		Permissions
	}

	Permissions struct {
		IncludedCountries map[string]bool
		IncludedProvinces map[string]map[string]bool
		ExcludedProvinces map[string]map[string]bool
		IncludedCities    map[string]map[string]map[string]bool
		ExcludedCities    map[string]map[string]map[string]bool
	}
)

func (c *Contract) AddIncludedRegion(regionString string) error {
	region, err := regions.GetRegionDetails(regionString)
	if err != nil {
		return err
	}

	switch region.Type {
	case regions.COUNTRY:
		c.IncludedCountries[region.CountryCode] = true
	case regions.PROVINCE:
		if c.IncludedProvinces[region.CountryCode] == nil {
			c.IncludedProvinces[region.CountryCode] = make(map[string]bool)
		}
		c.IncludedProvinces[region.CountryCode][region.ProvinceCode] = true
	case regions.CITY:
		if c.IncludedCities[region.CountryCode] == nil {
			c.IncludedCities[region.CountryCode] = make(map[string]map[string]bool)
		}
		if c.IncludedCities[region.CountryCode][region.ProvinceCode] == nil {
			c.IncludedCities[region.CountryCode][region.ProvinceCode] = make(map[string]bool)
		}
		c.IncludedCities[region.CountryCode][region.ProvinceCode][region.CityCode] = true
	}
	return nil
}

func (c *Contract) AddExcludedRegion(regionString string) error {
	region, err := regions.GetRegionDetails(regionString)
	if err != nil {
		return err
	}

	switch region.Type {
	case regions.COUNTRY: //its meaningless to exclude a country, as there is no world level inclusion to exclude from
	// 	c.IncludedCountries[region.CountryCode] = false
	case regions.PROVINCE:
		if c.ExcludedProvinces[region.CountryCode] == nil {
			c.ExcludedProvinces[region.CountryCode] = make(map[string]bool)
		}
		c.ExcludedProvinces[region.CountryCode][region.ProvinceCode] = true
	case regions.CITY:
		if c.ExcludedCities[region.CountryCode] == nil {
			c.ExcludedCities[region.CountryCode] = make(map[string]map[string]bool)
		}
		if c.ExcludedCities[region.CountryCode][region.ProvinceCode] == nil {
			c.ExcludedCities[region.CountryCode][region.ProvinceCode] = make(map[string]bool)
		}
		c.ExcludedCities[region.CountryCode][region.ProvinceCode][region.CityCode] = true
	}
	return nil
}
