package data

import (
	"challenge16/internal/regions"
	"challenge16/internal/response"
	"errors"
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
		Distributors map[string]distributorData
		mu           sync.RWMutex
	}

	distributorData struct {
		permissionDataGlobally
		parentDistributor *string
	}

	permissionDataGlobally map[string]permissionDataInCountry

	permissionDataInCountry struct {
		PermissionType string // "allow-all", "deny-all", "custom"
		Inclusions     map[string]permissionDataInProvince
		Exclusions     map[string]permissionDataInProvince
	}

	permissionDataInProvince struct {
		PermissionType string // "allow-all", "deny-all", "custom"
		Inclusions     map[string]bool
		Exclusions     map[string]bool
	}
)

func NewDataBank() DataBank {
	return DataBank{
		Distributors: make(map[string]distributorData),
		mu:           sync.RWMutex{},
	}
}

func (db *DataBank) MarkInclusion(distributor, regionString string) response.Response {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Distributors[distributor]; !ok {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, ErrDistributorNotFound)
	}

	countryCode, provinceCode, cityCode, regionType, err := db.getRegionDetails(regionString)
	if err != nil {
		return response.CreateError(404, REGION_NOT_FOUND, err)
	}

	if regionType == COUNTRY {
		db.Distributors[distributor].permissionDataGlobally[countryCode] = permissionDataInCountry{
			PermissionType: allowAll,
		}
		return successResponse
	}
	if _, ok := db.Distributors[distributor].permissionDataGlobally[countryCode]; !ok {
		db.Distributors[distributor].permissionDataGlobally[countryCode] = permissionDataInCountry{
			PermissionType: custom,
			Inclusions:     make(map[string]permissionDataInProvince),
		}
	}
	if regionType == PROVINCE {
		db.Distributors[distributor].permissionDataGlobally[countryCode].Inclusions[provinceCode] = permissionDataInProvince{
			PermissionType: allowAll,
		}
		return successResponse
	}

	if _, ok := db.Distributors[distributor].permissionDataGlobally[countryCode].Inclusions[provinceCode]; !ok {
		db.Distributors[distributor].permissionDataGlobally[countryCode].Inclusions[provinceCode] = permissionDataInProvince{
			PermissionType: custom,
			Inclusions:     make(map[string]bool),
		}
	}

	db.Distributors[distributor].permissionDataGlobally[countryCode].Inclusions[provinceCode].Inclusions[cityCode] = true
	return successResponse
}

func (db *DataBank) MarkExclusion(distributor, regionString string) response.Response {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Distributors[distributor]; !ok {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, ErrDistributorNotFound)
	}

	countryCode, provinceCode, cityCode, regionType, err := db.getRegionDetails(regionString)
	if err != nil {
		return response.CreateError(404, REGION_NOT_FOUND, err)
	}

	if regionType == COUNTRY {
		db.Distributors[distributor].permissionDataGlobally[countryCode] = permissionDataInCountry{
			PermissionType: denyAll,
		}
		return successResponse
	}
	if _, ok := db.Distributors[distributor].permissionDataGlobally[countryCode]; !ok {
		db.Distributors[distributor].permissionDataGlobally[countryCode] = permissionDataInCountry{
			PermissionType: custom,
			Inclusions:     make(map[string]permissionDataInProvince),
		}
	}
	if regionType == PROVINCE {
		db.Distributors[distributor].permissionDataGlobally[countryCode].Inclusions[provinceCode] = permissionDataInProvince{
			PermissionType: denyAll,
		}
		return successResponse
	}

	if _, ok := db.Distributors[distributor].permissionDataGlobally[countryCode].Inclusions[provinceCode]; !ok {
		db.Distributors[distributor].permissionDataGlobally[countryCode].Inclusions[provinceCode] = permissionDataInProvince{
			PermissionType: custom,
			Inclusions:     make(map[string]bool),
		}
	}

	db.Distributors[distributor].permissionDataGlobally[countryCode].Inclusions[provinceCode].Inclusions[cityCode] = false

	return successResponse
}

func (db *DataBank) getRegionDetails(regionString string) (countryCode, provinceCode, cityCode, regionType string, err error) {
	subStrings := strings.Split(regionString, "-") // Splitting the regionString by "-", this is the regionString I am assuming
	switch len(subStrings) {
	case 1:
		countryCode = subStrings[0]
		if !regions.CheckCountry(countryCode) {
			err = errors.New("country not found")
			return
		}
		regionType = COUNTRY
	case 2:
		countryCode = subStrings[1]
		provinceCode = subStrings[0]
		regionType = PROVINCE
		if !regions.CheckProvince(countryCode, provinceCode) {
			err = errors.New("country/province not found")
			return
		}
	default:
		countryCode = subStrings[2]
		provinceCode = subStrings[1]
		cityCode = subStrings[0]
		regionType = CITY
		if !regions.CheckCity(countryCode, provinceCode, cityCode) {
			err = errors.New("country/province/city not found")
			return
		}
	}
	return
}

func (db *DataBank) IsAllowed(distributor, regionString string) response.Response {
	countryCode, provinceCode, cityCode, regionType, err := db.getRegionDetails(regionString)
	if err != nil {
		return response.CreateError(404, REGION_NOT_FOUND, err)
	}
	isAllowed, err := db.isAllowedForTheDistributor(distributor, countryCode, provinceCode, cityCode, regionType)
	if err != nil {
		if err == ErrDistributorNotFound {
			return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, err)
		} else {
			return response.CreateError(500, INTERNAL_SERVER_ERROR, err)
		}
	}

	if !isAllowed {
		return response.CreateError(200, "DISTRIBUTION_NOT_ALLOWED", nil)
	} else {
		return response.CreateSuccess(200, "DISTRIBUTION_ALLOWED", nil)
	}
}

func (db *DataBank) isAllowedForTheDistributor(distributor, countryCode, provinceCode, cityCode, regionType string) (bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	permissionData, ok := db.Distributors[distributor]
	if !ok {
		return false, ErrDistributorNotFound
	}

	if permissionData.parentDistributor != nil {
		if allowed, err := db.isAllowedForTheDistributor(*permissionData.parentDistributor, countryCode, provinceCode, cityCode, regionType); err != nil {
			return false, err
		} else if !allowed {
			return false, nil
		}
	}

	if permissionDataInCountry, exists := permissionData.permissionDataGlobally[countryCode]; exists {
		switch permissionDataInCountry.PermissionType {
		case allowAll:
			return true, nil
		case denyAll:
			return false, nil
		default:
			if regionType == COUNTRY {
				return false, nil
			}
		}
		if permissionDataInProvince, exists := permissionDataInCountry.Inclusions[provinceCode]; exists {
			switch permissionDataInProvince.PermissionType {
			case allowAll:
				return true, nil
			case denyAll:
				return false, nil
			default:
				if regionType == PROVINCE {
					return false, nil
				}
				if permissionDataInProvince.Inclusions[cityCode] {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (db *DataBank) AddSubDistributor(subDistributor, parentDistributor string) response.Response {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Distributors[parentDistributor]; !ok {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, ErrDistributorNotFound)
	}
	if _, ok := db.Distributors[subDistributor]; !ok {
		db.Distributors[subDistributor] = distributorData{
			permissionDataGlobally: make(permissionDataGlobally),
			parentDistributor:      &parentDistributor,
		}
	}
	return createdResponse
}

func (db *DataBank) AddDistributor(distributor string) response.Response {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Distributors[distributor]; ok {
		return response.CreateError(400, "DISTRIBUTOR_EXISTS", ErrDistributorExists)
	}
	db.Distributors[distributor] = distributorData{
		permissionDataGlobally: make(permissionDataGlobally),
		parentDistributor:      nil,
	}
	return createdResponse
}

func (db *DataBank) RemoveDistributor(distributor string) response.Response {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Distributors[distributor]; !ok {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, ErrDistributorNotFound)
	}
	delete(db.Distributors, distributor)
	return successResponse
}
