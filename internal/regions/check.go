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

func checkCountry(countryCode string) bool {
	_, ok := Countries[countryCode]
	return ok
}

func checkProvince(countryCode, provinceCode string) bool {
	if checkCountry(countryCode) == false {
		return false
	}
	_, ok := Countries[countryCode].Provinces[provinceCode]
	return ok
}

func checkCity(countryCode, provinceCode, cityCode string) bool {
	if checkProvince(countryCode, provinceCode) == false {
		return false
	}
	_, ok := Countries[countryCode].Provinces[provinceCode].Cities[cityCode]
	return ok
}

func GetRegionDetails(regionString string) (countryCode, provinceCode, cityCode, regionType string, err error) {
	subStrings := strings.Split(regionString, "-") // Splitting the regionString by "-", this is the regionString I am assuming
	switch len(subStrings) {
	case 1:
		countryCode = subStrings[0]
		if !checkCountry(countryCode) {
			err = errors.New("country not found: " + countryCode)
			return
		}
		regionType = COUNTRY
	case 2:
		countryCode = subStrings[1]
		provinceCode = subStrings[0]
		regionType = PROVINCE
		if !checkProvince(countryCode, provinceCode) {
			err = errors.New("country/province not found: " + countryCode + "-" + provinceCode)
			return
		}
	default:
		countryCode = subStrings[2]
		provinceCode = subStrings[1]
		cityCode = subStrings[0]
		regionType = CITY
		if !checkCity(countryCode, provinceCode, cityCode) {
			err = errors.New("country/province/city not found: " + countryCode + "-" + provinceCode + "-" + cityCode)
			return
		}
	}
	return
}
