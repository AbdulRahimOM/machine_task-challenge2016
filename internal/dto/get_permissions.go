package dto

type GetPermissionsData struct {
	Distributor string   `json:"distributor"`
	Included    []string `json:"included"`
	Excluded    []string `json:"excluded"`
}
