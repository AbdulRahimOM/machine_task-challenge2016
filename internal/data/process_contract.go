package data

import (
	"challenge16/internal/dto"
	"challenge16/internal/response"
	"fmt"
)

func mergeMapIntoMap[A map[string]bool](to, from A) {
	for k, v := range from {
		to[k] = v
	}
}

// filterContractPermissionsBasedOnParentPermissions filters the contract permissions based on the parent permissions.
// It removes the regions that are not included in the parent permissions, but included in the contract permissions.
// It also removes the regions that are excluded in the parent permissions, but not excluded in the contract permissions.
// It also removes the regions that are excluded in the parent permissions, but included in the contract permissions.
// If parent is nil or not existing, it does nothing.
func (db *DataBank) filterContractPermissionsBasedOnParentPermissions(contract dto.Contract) {
	if contract.ParentDistributor == nil {
		return
	}

	parentPermission, ok := db.getDistributorPermissionCopy(*contract.ParentDistributor)
	if !ok {
		return
	}

	contractPermissions := contract.Permissions

	for country := range contractPermissions.IncludedCountries {
		if parentPermission.includedCountries[country] {
			continue
		} else {
			delete(contractPermissions.IncludedCountries, country)
			if _, exists := contractPermissions.IncludedProvinces[country]; !exists {
				contractPermissions.IncludedProvinces[country] = make(map[string]bool)
			}
			for parentProvince := range parentPermission.includedProvinces[country] {
				contractPermissions.IncludedProvinces[country][parentProvince] = true
			}
			if _, exists := contractPermissions.IncludedCities[country]; !exists {
				contractPermissions.IncludedCities[country] = make(map[string]map[string]bool)
				for parentProvince := range parentPermission.includedCities[country] {
					if _, exists := contractPermissions.IncludedCities[country][parentProvince]; !exists {
						contractPermissions.IncludedCities[country][parentProvince] = make(map[string]bool)
					}
					for parentCity := range parentPermission.includedCities[country][parentProvince] {
						contractPermissions.IncludedCities[country][parentProvince][parentCity] = true
					}
				}
			}
		}
	}

	for country := range contractPermissions.IncludedProvinces {
		if parentPermission.includedCountries[country] {
			continue
		}
		for province := range contractPermissions.IncludedProvinces[country] {
			if _, exists := parentPermission.includedProvinces[country]; exists && parentPermission.includedProvinces[country][province] {
				continue
			} else {
				delete(contractPermissions.IncludedProvinces[country], province)
				if _, exists := parentPermission.includedCities[country]; exists {
					if _, exists := contractPermissions.IncludedCities[country]; !exists {
						contractPermissions.IncludedCities[country] = make(map[string]map[string]bool)
					}
					if _, exists := contractPermissions.IncludedCities[country][province]; !exists {
						contractPermissions.IncludedCities[country][province] = make(map[string]bool)
					}
					for parentCity := range parentPermission.includedCities[country][province] {
						contractPermissions.IncludedCities[country][province][parentCity] = true
					}
				}
			}
		}
	}

	for country := range contractPermissions.IncludedCities {
		if parentPermission.includedCountries[country] {
			continue
		}
		for province := range contractPermissions.IncludedCities[country] {
			if _, exists := parentPermission.includedProvinces[country]; exists && parentPermission.includedProvinces[country][province] {
				continue
			}
			for city := range contractPermissions.IncludedCities[country][province] {
				if _, exists := parentPermission.includedCities[country]; exists {
					if _, exists := parentPermission.includedCities[country][province]; exists && parentPermission.includedCities[country][province][city] {
						continue
					} else {
						delete(contractPermissions.IncludedCities[country][province], city)
					}
				} else {
					delete(contractPermissions.IncludedCities[country][province], city)
				}
			}
		}
	}

	//merging applicable excluded regions
	//provincial level exclusion
	for country := range parentPermission.excludedProvinces {
		if contractPermissions.IncludedCountries[country] {
			if _, exists := contractPermissions.ExcludedProvinces[country]; !exists {
				contractPermissions.ExcludedProvinces[country] = make(map[string]bool)
			}
			mergeMapIntoMap(contractPermissions.ExcludedProvinces[country], parentPermission.excludedProvinces[country])
			continue
		}
		if _, exists := contractPermissions.IncludedProvinces[country]; exists {
			for parentExcludedProvince := range parentPermission.excludedProvinces[country] {
				//delete the province from included provinces if it is excluded in parent
				delete(contractPermissions.IncludedProvinces[country], parentExcludedProvince)

				//no need of city level exclusion for cities in this province as the province is excluded
				if _, exists := contractPermissions.ExcludedCities[country]; exists {
					delete(contractPermissions.ExcludedCities[country], parentExcludedProvince)
				}
			}
		}

		if _, exists := contractPermissions.IncludedCities[country]; exists {
			//delete the cities in province from included cities as the province is excluded in parent
			for parentExcludedProvince := range parentPermission.excludedProvinces[country] {
				delete(contractPermissions.IncludedCities[country], parentExcludedProvince)
			}
		}
	}

	//city level exclusion
	for country := range parentPermission.excludedCities {
		for province := range parentPermission.excludedCities[country] {
			for city := range parentPermission.excludedCities[country][province] {
				if contractPermissions.IncludedCountries[country] {
					if _, exists := contractPermissions.ExcludedProvinces[country]; !exists || !contractPermissions.ExcludedProvinces[country][province] {
						//if the province is not excluded in contract, but country is included, then exclude the city
						if _, exists := contractPermissions.ExcludedCities[country]; !exists {
							contractPermissions.ExcludedCities[country] = make(map[string]map[string]bool)
						}
						if _, exists := contractPermissions.ExcludedCities[country][province]; !exists {
							contractPermissions.ExcludedCities[country][province] = make(map[string]bool)
						}
						contractPermissions.ExcludedCities[country][province][city] = true
					}
				} else {
					if _, exists := contractPermissions.IncludedProvinces[country]; exists && contractPermissions.IncludedProvinces[country][province] {
						contractPermissions.ExcludedCities[country][province][city] = true
					} else {
						if _, exists := contractPermissions.IncludedCities[country]; exists {
							if _, exists := contractPermissions.IncludedCities[country][province]; exists {
								delete(contractPermissions.IncludedCities[country][province], city)
							}
						}
					}
				}
			}
		}
	}

	//city level exclusion
	for country := range parentPermission.excludedCities {
		for province := range parentPermission.excludedCities[country] {
			for city := range parentPermission.excludedCities[country][province] {
				if contractPermissions.IncludedCountries[country] {
					if _, exists := contractPermissions.ExcludedProvinces[country]; !exists || !contractPermissions.ExcludedProvinces[country][province] {
						//=> the province is not excluded in contract, but country is included. So, exclude the city(if not already excluded)

						if _, exists := contractPermissions.ExcludedCities[country]; !exists {
							contractPermissions.ExcludedCities[country] = make(map[string]map[string]bool)
						}
						if _, exists := contractPermissions.ExcludedCities[country][province]; !exists {
							contractPermissions.ExcludedCities[country][province] = make(map[string]bool)
						}
						contractPermissions.ExcludedCities[country][province][city] = true
					}
				} else {
					if _, exists := contractPermissions.IncludedProvinces[country]; exists && contractPermissions.IncludedProvinces[country][province] {
						//=> the province is included in contract. So, exclude the city(if not already excluded)
						contractPermissions.ExcludedCities[country][province][city] = true
					} else {
						//=> the province is not included in contract. So, there wont be any exclusions required for cities in this province.
						//=> But, as country and province are not included, these cities may be in included list. So, we need to remove the city from included cities(if it is included)
						if _, exists := contractPermissions.IncludedCities[country]; exists {
							if _, exists := contractPermissions.IncludedCities[country][province]; exists {
								delete(contractPermissions.IncludedCities[country][province], city)
							}
						}
					}
				}
			}
		}
	}

	contract.Permissions = contractPermissions
}

func validateContract(contract dto.Contract) error {
	//if a region is included, sub regions should only be of 'excluded' type
	for country := range contract.IncludedCountries {
		if _, exists := contract.IncludedProvinces[country]; exists && len(contract.IncludedProvinces[country]) > 0 {
			return fmt.Errorf("country %s is included, but provinces are also included. There should only be exclusions of sub-regions for an included region", country)
		}
		for province := range contract.IncludedCities[country] {
			if len(contract.IncludedCities[country][province]) > 0 {
				return fmt.Errorf("country %s is included, but cities in province %s are also included. There should only be exclusions of sub-regions for an included region", country, province)
			}
		}
	}

	for country := range contract.IncludedProvinces {
		for province := range contract.IncludedProvinces[country] {
			if _, exists := contract.IncludedCities[country]; exists && len(contract.IncludedCities[country][province]) > 0 {
				return fmt.Errorf("province %s in country %s is included, but cities are also included. There should only be exclusions of sub-regions for an included region", province, country)
			}

			//same province should not be included and excluded
			if _, exists := contract.ExcludedProvinces[country]; exists && contract.ExcludedProvinces[country][province] {
				return fmt.Errorf("province %s in country %s is included and excluded. It should be either included or excluded", province, country)
			}
		}
	}

	for country := range contract.IncludedCities {
		for province := range contract.IncludedCities[country] {
			for city := range contract.IncludedCities[country][province] {
				//same city should not be included and excluded
				if _, exists := contract.ExcludedCities[country]; exists {
					if _, exists := contract.ExcludedCities[country][province]; exists && contract.ExcludedCities[country][province][city] {
						return fmt.Errorf("city %s in province %s in country %s is included and excluded. It should be either included or excluded", city, province, country)
					}
				}
			}
		}
	}

	for country := range contract.ExcludedProvinces {
		for province := range contract.ExcludedProvinces[country] {
			if _, exists := contract.IncludedCities[country]; exists && len(contract.IncludedCities[country][province]) > 0 {
				return fmt.Errorf("province %s in country %s is excluded, but its cities are included. A region cannot be excluded while including its sub-regions", province, country)
			}

			// A province can be excluded only if its country is included; otherwise, it's meaningless.
			if !contract.IncludedCountries[country] {
				return fmt.Errorf("province %s in country %s is excluded, but the country is not included. A region can be excluded only if its parent is included.", province, country)  
			}

			// A province should not be included and excluded at the same time
			if _, exists := contract.IncludedProvinces[country]; exists && contract.IncludedProvinces[country][province] {
				return fmt.Errorf("province %s in country %s cannot be both included and excluded", province, country)  
			}
		}
	}

	for country := range contract.ExcludedCities {
		for province := range contract.ExcludedCities[country] {
			for city := range contract.ExcludedCities[country][province] {
				// A city can be excluded only when either its province is included or its country is included without excluding the province.
				if !contract.IncludedCountries[country] {
					if _, exists := contract.IncludedProvinces[country]; !exists || !contract.IncludedProvinces[country][province] {
						return fmt.Errorf("city %s in province %s in country %s is excluded, but the country is not included and the province is not included. A region cannot be excluded while its parent region is not included", city, province, country)
					}
				}

				// A city should not be included and excluded at the same time
				if _, exists := contract.IncludedCities[country]; exists {
					if _, exists := contract.IncludedCities[country][province]; exists && contract.IncludedCities[country][province][city] {
						return fmt.Errorf("city %s in province %s in country %s cannot be both included and excluded", city, province, country)
					}
				}
			}
		}
	}

	return nil
}

func (db *DataBank) applyContractOnDistributor(finalContract dto.Contract) {
	recipient := finalContract.ContractRecipient

	db.createDistributorIfNotExists(recipient)

	oldPermissionData, _ := db.getDistributorPermissionCopy(recipient)
	newPermissionData := oldPermissionData.copyPermissionData()

	//merge included countries
	mergeMapIntoMap(newPermissionData.includedCountries, finalContract.IncludedCountries)

	//merge included provinces
	for country, provinces := range finalContract.IncludedProvinces {
		if _, exists := newPermissionData.includedProvinces[country]; !exists {
			newPermissionData.includedProvinces[country] = make(map[string]bool)
		}
		mergeMapIntoMap(newPermissionData.includedProvinces[country], provinces)
	}

	//merge included cities
	for country := range finalContract.IncludedCities {
		if _, exists := newPermissionData.includedCities[country]; !exists {
			newPermissionData.includedCities[country] = make(map[string]map[string]bool)
		}
		for province, cities := range finalContract.IncludedCities[country] {
			if _, exists := newPermissionData.includedCities[country][province]; !exists {
				newPermissionData.includedCities[country][province] = make(map[string]bool)
			}
			mergeMapIntoMap(newPermissionData.includedCities[country][province], cities)
		}
	}

	finalExcludedProvinces := map[string]map[string]bool{}

	// finding exclusions that are excluded for one, but not included for the other
	for country := range finalContract.ExcludedProvinces {
		for province := range finalContract.ExcludedProvinces[country] {
			/*
			   possiblity of inclusions:
			   country included,province not excluded=>included province
			   country not included,province included=>included province
			*/
			isIncludedInOther := false
			if oldPermissionData.includedCountries[country] {
				if _, exists := oldPermissionData.excludedProvinces[country]; !exists || !oldPermissionData.excludedProvinces[country][province] {
					isIncludedInOther = true
				}
			} else {
				if _, exists := oldPermissionData.includedProvinces[country]; exists && oldPermissionData.includedProvinces[country][province] {
					isIncludedInOther = true
				}
			}
			if !isIncludedInOther {
				if _, exists := finalExcludedProvinces[country]; !exists {
					finalExcludedProvinces[country] = make(map[string]bool)
				}
				finalExcludedProvinces[country][province] = true
			}
		}
	}
	for country := range oldPermissionData.excludedProvinces {
		for province := range oldPermissionData.excludedProvinces[country] {
			isIncludedInOther := false
			if finalContract.IncludedCountries[country] {
				if _, exists := finalContract.ExcludedProvinces[country]; !exists || !finalContract.ExcludedProvinces[country][province] {
					isIncludedInOther = true
				}
			} else {
				if _, exists := finalContract.IncludedProvinces[country]; exists && finalContract.IncludedProvinces[country][province] {
					isIncludedInOther = true
				}
			}
			if !isIncludedInOther {
				if _, exists := finalExcludedProvinces[country]; !exists {
					finalExcludedProvinces[country] = make(map[string]bool)
				}
				finalExcludedProvinces[country][province] = true
			}
		}
	}

	finalExcludedCities := map[string]map[string]map[string]bool{}

	//merge commonly excluded cities
	/*
	   possiblity of inclusions:
	   country included, province not excluded and city not excluded => included city
	   country not included, province included and city not excluded => included city
	   country not included, province not included and city included => included city
	*/
	for country := range finalContract.ExcludedCities {
		for province := range finalContract.ExcludedCities[country] {
			for city := range finalContract.ExcludedCities[country][province] {
				isIncludedInOther := true
				if oldPermissionData.includedCountries[country] {
					if _, exists := oldPermissionData.excludedProvinces[country]; !exists || !oldPermissionData.excludedProvinces[country][province] {
						if _, exists := oldPermissionData.excludedCities[country]; exists {
							if _, exists := oldPermissionData.excludedCities[country][province]; exists && oldPermissionData.excludedCities[country][province][city] {
								isIncludedInOther = false
							}
						}
					} else {
						isIncludedInOther = false
					}
				} else {
					if _, exists := oldPermissionData.includedProvinces[country]; exists && oldPermissionData.includedProvinces[country][province] {
						if _, exists := oldPermissionData.excludedCities[country]; exists {
							if _, exists := oldPermissionData.excludedCities[country][province]; exists && oldPermissionData.excludedCities[country][province][city] {
								isIncludedInOther = false
							}
						}
					} else {
						if _, exists := oldPermissionData.includedCities[country]; exists {
							if _, exists := oldPermissionData.includedCities[country][province]; !exists || !oldPermissionData.includedCities[country][province][city] {
								isIncludedInOther = false
							}
						} else {
							isIncludedInOther = false
						}
					}
				}

				if !isIncludedInOther {
					if _, exists := finalExcludedCities[country]; !exists {
						finalExcludedCities[country] = make(map[string]map[string]bool)
					}
					if _, exists := finalExcludedCities[country][province]; !exists {
						finalExcludedCities[country][province] = make(map[string]bool)
					}
					finalExcludedCities[country][province][city] = true
				}
			}
		}
	}

	for country := range oldPermissionData.excludedCities {
		for province := range oldPermissionData.excludedCities[country] {
			for city := range oldPermissionData.excludedCities[country][province] {
				isIncludedInOther := true
				if finalContract.IncludedCountries[country] {
					if _, exists := finalContract.ExcludedProvinces[country]; !exists || !finalContract.ExcludedProvinces[country][province] {
						if _, exists := finalContract.ExcludedCities[country]; exists {
							if _, exists := finalContract.ExcludedCities[country][province]; exists && finalContract.ExcludedCities[country][province][city] {
								isIncludedInOther = false
							}
						}
					} else {
						isIncludedInOther = false
					}
				} else {
					if _, exists := finalContract.IncludedProvinces[country]; exists && finalContract.IncludedProvinces[country][province] {
						if _, exists := finalContract.ExcludedCities[country]; exists {
							if _, exists := finalContract.ExcludedCities[country][province]; exists && finalContract.ExcludedCities[country][province][city] {
								isIncludedInOther = false
							}
						}
					} else {
						if _, exists := finalContract.IncludedCities[country]; exists {
							if _, exists := finalContract.IncludedCities[country][province]; !exists || !finalContract.IncludedCities[country][province][city] {
								isIncludedInOther = false
							}
						} else {
							isIncludedInOther = false
						}
					}
				}

				if !isIncludedInOther {
					if _, exists := finalExcludedCities[country]; !exists {
						finalExcludedCities[country] = make(map[string]map[string]bool)
					}
					if _, exists := finalExcludedCities[country][province]; !exists {
						finalExcludedCities[country][province] = make(map[string]bool)
					}
					finalExcludedCities[country][province][city] = true
				}

			}
		}

	}

	for country := range finalContract.ExcludedCities {
		if _, exists := oldPermissionData.excludedCities[country]; !exists {
			continue
		}
		for province := range finalContract.ExcludedCities[country] {
			if _, exists := oldPermissionData.excludedCities[country][province]; !exists {
				continue
			}
			for city := range finalContract.ExcludedCities[country][province] {
				if oldPermissionData.excludedCities[country][province][city] {
					if _, exists := finalExcludedCities[country]; !exists {
						finalExcludedCities[country] = make(map[string]map[string]bool)
					}
					if _, exists := finalExcludedCities[country][province]; !exists {
						finalExcludedCities[country][province] = make(map[string]bool)
					}
					finalExcludedCities[country][province][city] = true
				}
			}
		}
	}

	//replace the existing data with the new data
	newPermissionData.excludedProvinces = finalExcludedProvinces
	newPermissionData.excludedCities = finalExcludedCities

	//replace the recipient's permission data with the new data
	db.mu.Lock()
	defer db.mu.Unlock()
	db.Distributors[recipient] = newPermissionData
}

func (db *DataBank) ApplyContract(contract dto.Contract) response.Response {

	err := validateContract(contract)
	if err != nil {
		return response.CreateError(400, "INVALID_CONTRACT", fmt.Errorf("invalid contract, err: %v", err))
	}

	if contract.ParentDistributor != nil {
		if !db.distributorExists(*contract.ParentDistributor) {
			return response.CreateError(404, "PARENT_DISTRIBUTOR_NOT_FOUND", fmt.Errorf("parent distributor %s not found", *contract.ParentDistributor))
		}
		db.filterContractPermissionsBasedOnParentPermissions(contract)
	}

	db.applyContractOnDistributor(contract)
	return successResponse
}

func (db *DataBank) createDistributorIfNotExists(distributor string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Distributors[distributor]; !ok {
		db.Distributors[distributor] = newPermissionData()
	}
}
