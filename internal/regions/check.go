package regions

import (
	"errors"
	"strings"
)

const (
	allowAll = "allow-all"
	denyAll  = "deny-all"
	custom   = "custom"

	COUNTRY  = "country"
	PROVINCE = "province"
	CITY     = "city"
)

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

type Region struct {
	CountryCode  string
	ProvinceCode string
	CityCode     string
	Type         string
}

func GetRegionDetails(regionString string) (Region, error) {
	var (
		region                                          Region
		err                                             error
		countryCode, provinceCode, cityCode, regionType string
	)
	subStrings := strings.Split(regionString, "-") // Splitting the regionString by "-", this is the regionString I am assuming
	switch len(subStrings) {
	case 1:
		countryCode = subStrings[0]
		regionType = COUNTRY
		if !CheckCountry(countryCode) {
			err = errors.New("country not found: " + countryCode)
		}
	case 2:
		countryCode = subStrings[1]
		provinceCode = subStrings[0]
		regionType = PROVINCE
		if !CheckProvince(countryCode, provinceCode) {
			err = errors.New("country/province not found: " + countryCode + "-" + provinceCode)
		}
	default:
		countryCode = subStrings[2]
		provinceCode = subStrings[1]
		cityCode = subStrings[0]
		regionType = CITY
		if !CheckCity(countryCode, provinceCode, cityCode) {
			err = errors.New("country/province/city not found: " + countryCode + "-" + provinceCode + "-" + cityCode)
		}
	}

	region = Region{
		CountryCode:  countryCode,
		ProvinceCode: provinceCode,
		CityCode:     cityCode,
		Type:         regionType,
	}
	return region, err
}
