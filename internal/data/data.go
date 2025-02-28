package data

import (
	"challenge16/internal/regions"
	"challenge16/internal/response"
	"errors"
	"fmt"
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
		Distributors map[string]permissionDataGlobal
		mu           sync.RWMutex
	}

	permissionDataGlobal map[string]permissionDataInCountry

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
		Distributors: make(map[string]permissionDataGlobal),
		mu:           sync.RWMutex{},
	}
}

func (db *DataBank) MarkInclusion(distributor, regionString string) response.Response {
	db.mu.RLock()
	if _, ok := db.Distributors[distributor]; !ok {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, ErrDistributorNotFound)
	}
	db.mu.RUnlock()

	countryCode, provinceCode, cityCode, regionType, err := regions.GetRegionDetails(regionString)
	if err != nil {
		return response.CreateError(404, REGION_NOT_FOUND, err)
	}

	db.markAsIncluded(distributor, countryCode, provinceCode, cityCode, regionType)
	return successResponse
}

func (db *DataBank) markAsIncluded(distributor, countryCode, provinceCode, cityCode, regionType string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.Distributors[distributor]; !ok {
		db.Distributors[distributor] = make(permissionDataGlobal)
	}

	if regionType == COUNTRY {
		db.Distributors[distributor][countryCode] = permissionDataInCountry{
			PermissionType: allowAll,
			Inclusions:     make(map[string]permissionDataInProvince),
			Exclusions:     make(map[string]permissionDataInProvince),
		}
		return
	}
	if _, ok := db.Distributors[distributor][countryCode]; !ok {
		db.Distributors[distributor][countryCode] = permissionDataInCountry{
			PermissionType: custom,
			Inclusions:     make(map[string]permissionDataInProvince),
			Exclusions:     make(map[string]permissionDataInProvince),
		}
	}
	if regionType == PROVINCE {
		db.Distributors[distributor][countryCode].Inclusions[provinceCode] = permissionDataInProvince{
			PermissionType: allowAll,
			Inclusions:     make(map[string]bool),
			Exclusions:     make(map[string]bool),
		}
		return
	}

	if _, ok := db.Distributors[distributor][countryCode].Inclusions[provinceCode]; !ok {
		db.Distributors[distributor][countryCode].Inclusions[provinceCode] = permissionDataInProvince{
			PermissionType: custom,
			Inclusions:     make(map[string]bool),
			Exclusions:     make(map[string]bool),
		}
	}

	db.Distributors[distributor][countryCode].Inclusions[provinceCode].Inclusions[cityCode] = true
	return
}

func (db *DataBank) MarkExclusion(distributor, regionString string) response.Response {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Distributors[distributor]; !ok {
		return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, ErrDistributorNotFound)
	}

	countryCode, provinceCode, cityCode, regionType, err := regions.GetRegionDetails(regionString)
	if err != nil {
		return response.CreateError(404, REGION_NOT_FOUND, err)
	}

	db.markAsExcluded(distributor, countryCode, provinceCode, cityCode, regionType)
	return successResponse
}

func (db *DataBank) markAsExcluded(distributor, countryCode, provinceCode, cityCode, regionType string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if regionType == COUNTRY {
		db.Distributors[distributor][countryCode] = permissionDataInCountry{
			PermissionType: denyAll,
			Inclusions:     make(map[string]permissionDataInProvince),
			Exclusions:     make(map[string]permissionDataInProvince),
		}
		return
	}
	if _, ok := db.Distributors[distributor][countryCode]; !ok {
		db.Distributors[distributor][countryCode] = permissionDataInCountry{
			PermissionType: custom,
			Inclusions:     make(map[string]permissionDataInProvince),
			Exclusions:     make(map[string]permissionDataInProvince),
		}
	}
	if regionType == PROVINCE {

		db.Distributors[distributor][countryCode].Exclusions[provinceCode] = permissionDataInProvince{
			PermissionType: denyAll,
			Inclusions:     make(map[string]bool),
			Exclusions:     make(map[string]bool),
		}
		return
	}

	if _, ok := db.Distributors[distributor][countryCode].Exclusions[provinceCode]; !ok {
		db.Distributors[distributor][countryCode].Exclusions[provinceCode] = permissionDataInProvince{
			PermissionType: custom,
			Inclusions:     make(map[string]bool),
			Exclusions:     make(map[string]bool),
		}
	}

	db.Distributors[distributor][countryCode].Exclusions[provinceCode].Exclusions[cityCode] = false
	return
}

func (db *DataBank) AddDistributor(distributor string) response.Response {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Distributors[distributor]; ok {
		return response.CreateError(400, "DISTRIBUTOR_EXISTS", ErrDistributorExists)
	}
	db.Distributors[distributor] = make(permissionDataGlobal)
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

func (db *DataBank) ApplyContract(distributorHeirarchy, includeRegions, excludeRegions []string) response.Response {
	if len(distributorHeirarchy) > 1 {
		//ensure that they have required permission
		for i := 1; i < len(distributorHeirarchy); i++ { //skip the first distributor as it may be a new distributor
			if _, ok := db.Distributors[distributorHeirarchy[i]]; !ok {
				return response.CreateError(404, DISTRIBUTOR_NOT_FOUND, fmt.Errorf("parent distributor %s not found", distributorHeirarchy[i]))
			}
		}

		for _, region := range includeRegions {
			countryCode, provinceCode, cityCode, regionType, err := regions.GetRegionDetails(region)
			if err != nil {
				return response.CreateError(404, REGION_NOT_FOUND, err)
			}

			//check if the region is allowed for the immediate parent distributor
			isAllowedForImmediateParent := db.isAllowedForTheDistributor(distributorHeirarchy[1], countryCode, provinceCode, cityCode, regionType)
			if !isAllowedForImmediateParent {
				return response.CreateError(200, "DISTRIBUTION_NOT_ALLOWED", fmt.Errorf("distribution not allowed for the immediate parent(%s) in region %s which is mentioned in 'INCLUDE'", distributorHeirarchy[1], region))
			}
		}
	}

	//validate the regions mentioned in the contract
	for _, region := range includeRegions {
		_, _, _, _, err := regions.GetRegionDetails(region)
		if err != nil {
			return response.CreateError(404, REGION_NOT_FOUND, err)
		}
	}

	for _, region := range excludeRegions {
		_, _, _, _, err := regions.GetRegionDetails(region)
		if err != nil {
			return response.CreateError(404, REGION_NOT_FOUND, err)
		}
	}

	if _, ok := db.Distributors[distributorHeirarchy[0]]; !ok {
		db.Distributors[distributorHeirarchy[0]] = make(permissionDataGlobal)
	}

	//apply the contract
	for _, includeRegion := range includeRegions {
		countryCode, provinceCode, cityCode, regionType, _ := regions.GetRegionDetails(includeRegion)
		db.markAsIncluded(distributorHeirarchy[0], countryCode, provinceCode, cityCode, regionType)
	}

	for _, excludeRegion := range excludeRegions {
		countryCode, provinceCode, cityCode, regionType, _ := regions.GetRegionDetails(excludeRegion)
		if !db.isAllowedForTheDistributor(distributorHeirarchy[0], countryCode, provinceCode, cityCode, regionType) {
			continue //if the region is already in allow list(possibly by other contracts), then no need to exclude it
		} else {
			db.markAsExcluded(distributorHeirarchy[0], countryCode, provinceCode, cityCode, regionType)
		}
	}

	return successResponse
}

func (db *DataBank) isAllowedForTheDistributor(distributor, countryCode, provinceCode, cityCode, regionType string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	permissionDataGlobally, ok := db.Distributors[distributor]
	if !ok {
		return false
	}

	if permissionDataInCountry, exists := permissionDataGlobally[countryCode]; exists {
		switch permissionDataInCountry.PermissionType {
		case allowAll:
			return true
		case denyAll:
			return false
		default:
			if regionType == COUNTRY {
				return false
			}
		}
		if permissionDataInProvince, exists := permissionDataInCountry.Inclusions[provinceCode]; exists {
			switch permissionDataInProvince.PermissionType {
			case allowAll:
				return true
			case denyAll:
				return false
			default:
				if regionType == PROVINCE {
					return false
				}
				if permissionDataInProvince.Inclusions[cityCode] {
					return true
				}
			}
		}
	}
	return false
}
