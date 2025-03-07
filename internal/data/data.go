package data

import (
	"challenge16/internal/dto"
	"challenge16/internal/regions"
	"challenge16/internal/response"
	"errors"
	"fmt"
	"strings"
	"sync"
)

const (
	allowAll = "allow-all"
	denyAll  = "deny-all"
	custom   = "custom"

	COUNTRY  = "country"
	PROVINCE = "province"
	CITY     = "city"

	DISTRIBUTOR_NOT_FOUND = "DISTRIBUTOR_NOT_FOUND"
	REGION_NOT_FOUND      = "REGION_NOT_FOUND"
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
)

var (
	ErrDistributorNotFound = errors.New("distributor not found")
	ErrDistributorExists   = errors.New("distributor already exists")

	successResponse = response.CreateSuccess(200, "SUCCESS", nil)
	createdResponse = response.CreateSuccess(201, "CREATED", nil)
)

type (
	DataBank struct {
		Distributors map[string]permissionData
		mu           sync.RWMutex
	}
)

func NewDataBank() DataBank {
	return DataBank{
		Distributors: make(map[string]permissionData),
		mu:           sync.RWMutex{},
	}
}

func newPermissionData() permissionData {
	return permissionData{
		includedCountries: make(map[string]bool),
		includedProvinces: make(map[string]map[string]bool),
		excludedProvinces: make(map[string]map[string]bool),
		includedCities:    make(map[string]map[string]map[string]bool),
		excludedCities:    make(map[string]map[string]map[string]bool),
	}
}

func (db *DataBank) MarkInclusion(distributor, regionString string) response.Response {
	if !db.distributorExists(distributor) {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, fmt.Errorf("distributor %s not found", distributor))
	}

	region, err := regions.GetRegionDetails(regionString)
	if err != nil {
		return response.CreateError(404, REGION_NOT_FOUND, err)
	}

	db.markAsIncluded(distributor, region)
	return successResponse
}

func (db *DataBank) markAsIncluded(distributor string, region regions.Region) {
	db.createDistributorIfNotExists(distributor)

	countryCode, provinceCode, cityCode := region.CountryCode, region.ProvinceCode, region.CityCode

	db.mu.Lock()
	defer db.mu.Unlock()

	switch region.Type {
	case COUNTRY:
		db.Distributors[distributor].includedCountries[countryCode] = true
		delete(db.Distributors[distributor].includedProvinces, countryCode)
		delete(db.Distributors[distributor].excludedProvinces, countryCode)
		delete(db.Distributors[distributor].includedCities, countryCode)
		delete(db.Distributors[distributor].excludedCities, countryCode)
	case PROVINCE:
		if db.Distributors[distributor].includedCountries[countryCode] {
			if _, exists := db.Distributors[distributor].excludedProvinces[countryCode]; exists {
				if db.Distributors[distributor].excludedProvinces[countryCode][provinceCode] {
					delete(db.Distributors[distributor].excludedProvinces[countryCode], provinceCode)

					if _, exists := db.Distributors[distributor].includedCities[countryCode]; exists {
						delete(db.Distributors[distributor].includedCities[countryCode], provinceCode)
					}
				}
			}
			if _, exists := db.Distributors[distributor].excludedCities[countryCode]; exists {
				delete(db.Distributors[distributor].excludedCities[countryCode], provinceCode)
			}
			//no need to check in includedProvinces, because country is already included. so there will be only exceptions
		} else {
			//no need to check in excludedProvinces, because country is not included to have exceptions
			if _, exists := db.Distributors[distributor].includedProvinces[countryCode]; exists {
				if db.Distributors[distributor].includedProvinces[countryCode][provinceCode] {

					//excudedCities are no longer excluded as the province is now included
					if _, exists := db.Distributors[distributor].excludedCities[countryCode]; exists {
						delete(db.Distributors[distributor].excludedCities[countryCode], provinceCode)
					}
				} else {
					db.Distributors[distributor].includedProvinces[countryCode][provinceCode] = true

					//deleting includedCities in this province as the province itself is now included as a whole
					if _, exists := db.Distributors[distributor].includedCities[countryCode]; exists {
						delete(db.Distributors[distributor].includedCities[countryCode], provinceCode)
					}
				}
			} else {
				db.Distributors[distributor].includedProvinces[countryCode] = map[string]bool{provinceCode: true}

				//deleting includedCities in this province as the province itself is now included as a whole
				if _, exists := db.Distributors[distributor].includedCities[countryCode]; exists {
					delete(db.Distributors[distributor].includedCities[countryCode], provinceCode)
				}
				//as country or province were not there as 'included', so no need to check in excludedCities either
			}
		}
	case CITY:
		if db.Distributors[distributor].includedCountries[countryCode] {
			//there is no need to check in includedProvinces, because country is already included. so there will be only exceptions

			//checking in excludedProvinces:
			if _, exists := db.Distributors[distributor].excludedProvinces[countryCode]; exists {
				if db.Distributors[distributor].excludedProvinces[countryCode][provinceCode] {

					// if the province is excluded, then the city should be made to be included
					if _, exists := db.Distributors[distributor].includedCities[countryCode]; exists {
						if _, exists := db.Distributors[distributor].includedCities[countryCode][provinceCode]; !exists {
							db.Distributors[distributor].includedCities[countryCode][provinceCode] = map[string]bool{cityCode: true}
						}
					} else {
						db.Distributors[distributor].includedCities[countryCode] = map[string]map[string]bool{provinceCode: {cityCode: true}}
					}
				}
			}

			//checking in excludedCities:
			if _, exists := db.Distributors[distributor].excludedCities[countryCode]; exists {
				if _, exists := db.Distributors[distributor].excludedCities[countryCode][provinceCode]; exists {
					delete(db.Distributors[distributor].excludedCities[countryCode][provinceCode], cityCode)
				}
			}
		} else {
			//as not even the country is included, we just have to check in includedProvinces+excludedCities and includedCities

			//checking in includedProvinces:
			if _, exists := db.Distributors[distributor].includedProvinces[countryCode]; exists {
				if db.Distributors[distributor].includedProvinces[countryCode][provinceCode] {

					//ensure that city is not excluded
					if _, exists := db.Distributors[distributor].excludedCities[countryCode]; exists {
						if _, exists := db.Distributors[distributor].excludedCities[countryCode][provinceCode]; exists {
							delete(db.Distributors[distributor].excludedCities[countryCode][provinceCode], cityCode)
						}
					}

				} else {
					//if the province is not included, then the city should be included
					if _, exists := db.Distributors[distributor].includedCities[countryCode]; exists {
						if _, exists := db.Distributors[distributor].includedCities[countryCode][provinceCode]; !exists {
							db.Distributors[distributor].includedCities[countryCode][provinceCode] = map[string]bool{cityCode: true}
						} else {
							db.Distributors[distributor].includedCities[countryCode][provinceCode][cityCode] = true
						}
					} else {
						db.Distributors[distributor].includedCities[countryCode] = map[string]map[string]bool{provinceCode: {cityCode: true}}
					}
				}
			} else {

				//if the province is not included along with country not being included, then the city should be included
				if _, exists := db.Distributors[distributor].includedCities[countryCode]; exists {
					if _, exists := db.Distributors[distributor].includedCities[countryCode][provinceCode]; !exists {
						db.Distributors[distributor].includedCities[countryCode][provinceCode] = map[string]bool{cityCode: true}
					} else {
						db.Distributors[distributor].includedCities[countryCode][provinceCode][cityCode] = true
					}
				} else {
					db.Distributors[distributor].includedCities[countryCode] = map[string]map[string]bool{provinceCode: {cityCode: true}}
				}

				//as country or province were not there as 'included', so need to check in excludedCities
				if _, exists := db.Distributors[distributor].excludedCities[countryCode]; exists {
					if _, exists := db.Distributors[distributor].excludedCities[countryCode][provinceCode]; exists {
						delete(db.Distributors[distributor].excludedCities[countryCode][provinceCode], cityCode)
					}
				}

			}
		}
	}
}

func (db *DataBank) MarkExclusion(distributor, regionString string) response.Response {
	if !db.distributorExists(distributor) {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, fmt.Errorf("distributor %s not found", distributor))
	}

	region, err := regions.GetRegionDetails(regionString)
	if err != nil {
		return response.CreateError(404, REGION_NOT_FOUND, err)
	}

	db.markAsExcluded(distributor, region)
	return successResponse
}

func (db *DataBank) markAsExcluded(distributor string, region regions.Region) {
	db.createDistributorIfNotExists(distributor)

	countryCode, provinceCode, cityCode := region.CountryCode, region.ProvinceCode, region.CityCode

	db.mu.Lock()
	defer db.mu.Unlock()

	switch region.Type {
	case COUNTRY:
		delete(db.Distributors[distributor].includedCountries, countryCode)
		delete(db.Distributors[distributor].includedProvinces, countryCode)
		delete(db.Distributors[distributor].includedCities, countryCode)
		delete(db.Distributors[distributor].excludedProvinces, countryCode)
		delete(db.Distributors[distributor].excludedCities, countryCode)

	case PROVINCE:
		if db.Distributors[distributor].includedCountries[countryCode] {
			//if the country is included, then the province should be excluded
			if _, exists := db.Distributors[distributor].excludedProvinces[countryCode]; exists {
				db.Distributors[distributor].excludedProvinces[countryCode][provinceCode] = true
			} else {
				db.Distributors[distributor].excludedProvinces[countryCode] = map[string]bool{provinceCode: true}
			}

			//as the province as a whole is excluded, then the cities in the province need not be in excluded list
			if _, exists := db.Distributors[distributor].excludedCities[countryCode]; exists {
				delete(db.Distributors[distributor].excludedCities[countryCode], provinceCode)
			}
		} else {
			//if the country is not included, then the province should not be in included list
			if _, exists := db.Distributors[distributor].includedProvinces[countryCode]; exists {
				if db.Distributors[distributor].includedProvinces[countryCode][provinceCode] {
					delete(db.Distributors[distributor].includedProvinces[countryCode], provinceCode)

					//as the province is excluded, then the cities in the province need not be in excluded list
					if _, exists := db.Distributors[distributor].excludedCities[countryCode]; exists {
						delete(db.Distributors[distributor].excludedCities[countryCode], provinceCode)
					}
				}
			}

			//ensure that no city in the province is in included list
			if _, exists := db.Distributors[distributor].includedCities[countryCode]; exists {
				delete(db.Distributors[distributor].includedCities[countryCode], provinceCode)
			}

		}

	case CITY:
		if db.Distributors[distributor].includedCountries[countryCode] {
			//check if the province is excluded
			if _, exists := db.Distributors[distributor].excludedProvinces[countryCode]; exists && db.Distributors[distributor].excludedProvinces[countryCode][provinceCode] {
				//if the province is excluded, then the city should not be in included list
				if _, exists := db.Distributors[distributor].includedCities[countryCode]; exists {
					if _, exists := db.Distributors[distributor].includedCities[countryCode][provinceCode]; exists {
						delete(db.Distributors[distributor].includedCities[countryCode][provinceCode], cityCode)
					}
				}
			} else { //if the province is not excluded, then the city should be in excluded list
				if _, exists := db.Distributors[distributor].excludedCities[countryCode]; exists {
					if _, exists := db.Distributors[distributor].excludedCities[countryCode][provinceCode]; exists {
						db.Distributors[distributor].excludedCities[countryCode][provinceCode][cityCode] = true
					} else {
						db.Distributors[distributor].excludedCities[countryCode][provinceCode] = map[string]bool{cityCode: true}
					}
				} else {
					db.Distributors[distributor].excludedCities[countryCode] = map[string]map[string]bool{provinceCode: {cityCode: true}}
				}
			}
		} else { //if the country is not included, then either (the city should be in excluded list) or (the province should be in excluded list with no exception for the city)
			if _, exists := db.Distributors[distributor].excludedProvinces[countryCode]; exists && db.Distributors[distributor].excludedProvinces[countryCode][provinceCode] {
				if _, exists := db.Distributors[distributor].excludedCities[countryCode]; exists {
					if _, exists := db.Distributors[distributor].excludedCities[countryCode][provinceCode]; exists {
						db.Distributors[distributor].excludedCities[countryCode][provinceCode][cityCode] = true
					} else {
						db.Distributors[distributor].excludedCities[countryCode][provinceCode] = map[string]bool{cityCode: true}
					}
				} else {
					db.Distributors[distributor].excludedCities[countryCode] = map[string]map[string]bool{provinceCode: {cityCode: true}}
				}
			} else { //if the province is not excluded, then the city should be in excluded list
				if _, exists := db.Distributors[distributor].excludedCities[countryCode]; exists {
					if _, exists := db.Distributors[distributor].excludedCities[countryCode][provinceCode]; exists {
						db.Distributors[distributor].excludedCities[countryCode][provinceCode][cityCode] = true
					} else {
						db.Distributors[distributor].excludedCities[countryCode][provinceCode] = map[string]bool{cityCode: true}
					}
				} else {
					db.Distributors[distributor].excludedCities[countryCode] = map[string]map[string]bool{provinceCode: {cityCode: true}}
				}
			}

		}

	}

}

func (db *DataBank) AddDistributor(distributor string) response.Response {
	if db.distributorExists(distributor) {
		return response.CreateError(400, "DISTRIBUTOR_EXISTS", ErrDistributorExists)
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Distributors[distributor] = newPermissionData()
	return createdResponse
}

func (db *DataBank) RemoveDistributor(distributor string) response.Response {
	if !db.distributorExists(distributor) {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, ErrDistributorNotFound)
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.Distributors, distributor)

	return successResponse
}

func (db *DataBank) isAllowedForTheDistributor(distributor string, region regions.Region) (bool, string) {
	const (
		PARTIALLY_ALLOWED = "PARTIALLY_ALLOWED"
		FULLY_ALLOWED     = "FULLY_ALLOWED"
		FULLY_DENIED      = "FULLY_DENIED"
	)
	db.mu.RLock()
	defer db.mu.RUnlock()
	permissionData, ok := db.Distributors[distributor]
	if !ok {
		return false, ""
	}

	countryCode, provinceCode, cityCode, regionType := region.CountryCode, region.ProvinceCode, region.CityCode, region.Type

	switch regionType {
	case COUNTRY:
		if permissionData.includedCountries[countryCode] {
			if len(permissionData.excludedProvinces[countryCode]) > 0 {
				return false, PARTIALLY_ALLOWED
			}
			if len(permissionData.excludedCities[countryCode]) > 0 {
				return false, PARTIALLY_ALLOWED
			}
			return true, FULLY_ALLOWED
		} else {
			if len(permissionData.includedProvinces[countryCode]) > 0 {
				return false, PARTIALLY_ALLOWED
			}
			for provinceCode := range permissionData.includedCities[countryCode] {
				if len(permissionData.includedCities[countryCode][provinceCode]) > 0 {
					return false, PARTIALLY_ALLOWED
				}
			}
			return false, FULLY_DENIED
		}
	case PROVINCE:
		if permissionData.includedCountries[countryCode] {
			if _, exists := permissionData.excludedProvinces[countryCode]; exists && permissionData.excludedProvinces[countryCode][provinceCode] {
				return false, FULLY_DENIED
			}
			if _, exists := permissionData.excludedCities[countryCode]; exists && len(permissionData.excludedCities[countryCode][provinceCode]) > 0 {
				return false, PARTIALLY_ALLOWED
			}
			return true, FULLY_ALLOWED
		} else {
			if _, exists := permissionData.includedProvinces[countryCode]; exists && permissionData.includedProvinces[countryCode][provinceCode] {
				if _, exists := permissionData.excludedCities[countryCode]; exists && len(permissionData.excludedCities[countryCode][provinceCode]) > 0 {
					return false, PARTIALLY_ALLOWED
				}
				return true, FULLY_ALLOWED
			} else {
				if _, exists := permissionData.includedCities[countryCode]; exists {
					if _, exists := permissionData.includedCities[countryCode][provinceCode]; exists && len(permissionData.includedCities[countryCode][provinceCode]) > 0 {
						return false, PARTIALLY_ALLOWED
					}
				}
			}
		}
	case CITY:
		if permissionData.includedCountries[countryCode] {
			if _, exists := permissionData.excludedProvinces[countryCode]; exists && permissionData.excludedProvinces[countryCode][provinceCode] {
				return false, FULLY_DENIED
			}
			if _, exists := permissionData.excludedCities[countryCode]; exists {
				if _, exists := permissionData.excludedCities[countryCode][provinceCode]; exists && permissionData.excludedCities[countryCode][provinceCode][cityCode] {
					return false, FULLY_DENIED
				}
				return true, FULLY_ALLOWED
			} else {
				if _, exists := permissionData.includedProvinces[countryCode]; exists && permissionData.includedProvinces[countryCode][provinceCode] {
					if _, exists := permissionData.excludedCities[countryCode]; exists {
						if _, exists := permissionData.excludedCities[countryCode][provinceCode]; exists && permissionData.excludedCities[countryCode][provinceCode][cityCode] {
							return false, FULLY_DENIED
						}
					}
					return true, FULLY_ALLOWED
				}
				if _, exists := permissionData.includedCities[countryCode]; exists {
					if _, exists := permissionData.includedCities[countryCode][provinceCode]; exists && permissionData.includedCities[countryCode][provinceCode][cityCode] {
						return true, FULLY_ALLOWED
					}
				}
				return false, FULLY_DENIED
			}
		}
	}

	return false, "UNKNOWN" //this should never happen, as the region type is already validated
}

func (db *DataBank) GetDistributors() response.Response {
	db.mu.RLock()
	defer db.mu.RUnlock()
	distributors := make([]string, 0, len(db.Distributors))
	for distributor := range db.Distributors {
		distributors = append(distributors, distributor)
	}
	return response.CreateSuccess(200, "SUCCESS", map[string]interface{}{
		"distributors": distributors,
	})
}

func (db *DataBank) CheckIfDistributionIsAllowed(distributor, regionString string) response.Response {
	region, err := regions.GetRegionDetails(regionString)
	if err != nil {
		return response.CreateError(404, REGION_NOT_FOUND, err)
	}

	if !db.distributorExists(distributor) {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, fmt.Errorf("distributor %s not found", distributor))
	}

	_, status := db.isAllowedForTheDistributor(distributor, region)
	return response.CreateSuccess(200, status, nil)
}

func (db *DataBank) distributorExists(distributor string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if _, ok := db.Distributors[distributor]; ok {
		return true
	}
	return false
}

// func (db *DataBank) getParentRegions(distributor string) ([]regions.Region, []regions.Region) {
// 	return nil, nil
// }

func (db *DataBank) getDistributorPermissionCopy(distributor string) (permissionData, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if permissionData, ok := db.Distributors[distributor]; ok {
		return permissionData.copyPermissionData(), true
	}
	return permissionData{}, false
}

func (db *DataBank) GetDistributorPermissionsAsText(distributor string) string {
	permissionData, ok := db.getDistributorPermissionCopy(distributor)
	if !ok {
		return "Distributor not found"
	}

	builder := new(strings.Builder)
	builder.WriteString("Permissions for " + distributor)

	for country := range permissionData.includedCountries {
		builder.WriteString("\nINCLUDE: " + country)
	}
	for country := range permissionData.includedProvinces {
		for province := range permissionData.includedProvinces[country] {
			builder.WriteString("\nINCLUDE: " + province + "-" + country)
		}
	}
	for country := range permissionData.includedCities {
		for province := range permissionData.includedCities[country] {
			for city := range permissionData.includedCities[country][province] {
				builder.WriteString("\nINCLUDE: " + city + "-" + province + "-" + country)
			}
		}
	}

	for country := range permissionData.excludedProvinces {
		for province := range permissionData.excludedProvinces[country] {
			builder.WriteString("\nEXCLUDE: " + province + "-" + country)
		}
	}
	for country := range permissionData.excludedCities {
		for province := range permissionData.excludedCities[country] {
			for city := range permissionData.excludedCities[country][province] {
				builder.WriteString("\nEXCLUDE: " + city + "-" + province + "-" + country)
			}
		}
	}

	return builder.String()
}

func (db *DataBank) GetDistributorPermissionAsJSON(distributor string) response.Response {
	permissionData, ok := db.getDistributorPermissionCopy(distributor)
	if !ok {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, fmt.Errorf("distributor %s not found", distributor))
	}

	inclusions := make([]string, 0, len(permissionData.includedCountries)+len(permissionData.includedProvinces)+len(permissionData.includedCities))
	exclusions := make([]string, 0, len(permissionData.excludedProvinces)+len(permissionData.excludedCities))

	for country := range permissionData.includedCountries {
		inclusions = append(inclusions, country)
	}
	for country := range permissionData.includedProvinces {
		for province := range permissionData.includedProvinces[country] {
			inclusions = append(inclusions, province+"-"+country)
		}
	}
	for country := range permissionData.includedCities {
		for province := range permissionData.includedCities[country] {
			for city := range permissionData.includedCities[country][province] {
				inclusions = append(inclusions, city+"-"+province+"-"+country)
			}
		}
	}

	for country := range permissionData.excludedProvinces {
		for province := range permissionData.excludedProvinces[country] {
			exclusions = append(exclusions, province+"-"+country)
		}
	}

	for country := range permissionData.excludedCities {
		for province := range permissionData.excludedCities[country] {
			for city := range permissionData.excludedCities[country][province] {
				exclusions = append(exclusions, city+"-"+province+"-"+country)
			}
		}
	}

	data:=dto.GetPermissionsData{
		Distributor: distributor,
		Included:    inclusions,
		Excluded:    exclusions,
	}

	return response.CreateSuccess(200, "SUCCESS", data)
}
