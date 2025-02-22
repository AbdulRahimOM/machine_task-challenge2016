package handler

import "challenge16/internal/data"

const (
	URL_PARAM_MISSING = "URL_PARAM_MISSING"
)

type handler struct {
	databank data.DataBank
}

func NewHandler() *handler {
	return &handler{
		databank: data.NewDataBank(),
	}
}
