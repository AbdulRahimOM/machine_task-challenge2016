package response

type Response struct {
	HttpStatusCode int         `json:"-"`
	Status         bool        `json:"status"`
	ResponseCode   string      `json:"resp_code"`
	Error          error       `json:"-"` //will be marshalled to string when WriteToJSON is called
	Data           interface{} `json:"data,omitempty"`
}

type ValidationErrorResponse struct {
	Status       bool           `json:"status"`
	ResponseCode string         `json:"resp_code"`
	Errors       []InvalidField `json:"errors"`
}

type InvalidField struct {
	FailedField string      `json:"field"`
	Tag         string      `json:"tag"`
	Value       interface{} `json:"value"`
}
