package response

import (
	"fmt"
	"net/http"
)

const (
	INVALID_URL_PARAM = "INVALID_URL_PARAM"
)

func CreateError(statusCode int, respcode string, err error) Response {
	return Response{
		HttpStatusCode: statusCode,
		Status:         false,
		ResponseCode:   respcode,
		Error:          err,
	}
}

func CreateSuccess(statusCode int, respcode string, data interface{}) Response {
	return Response{
		HttpStatusCode: statusCode,
		Status:         true,
		ResponseCode:   respcode,
		Data:           data,
	}
}

func InvalidURLParamResponse(param string, err error) Response {
	return CreateError(http.StatusBadRequest, INVALID_URL_PARAM, fmt.Errorf("error parsing %v from url: %w", param, err))
}
