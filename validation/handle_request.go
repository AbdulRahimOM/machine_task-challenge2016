package validation

import (
	"challenge16/internal/response"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const (
	bindingErrCode    = "BINDING_ERROR"
	validationErrCode = "VALIDATION_ERROR"

	fieldTag_JSON  = "json"
	fieldTag_Query = "query"
)

/*
BindAndValidateRequest binds and validates the request.
Req should be a pointer to the request struct.
*/
func BindAndValidateJSONRequest(c *fiber.Ctx, req interface{}) (bool, error) {
	if err := c.BodyParser(req); err != nil {
		return false, response.Response{
			HttpStatusCode: 400,
			Status:         false,
			ResponseCode:   bindingErrCode,
			Error:          fmt.Errorf("error parsing request:%w", err),
		}.WriteToJSON(c)
	}

	if ok, errResponse := validateRequestInDetailBasedOnFieldTag(c, req, fieldTag_JSON); !ok {
		return false, errResponse
	}

	log.Debug("req after validation:", req) //alter later if need to hide sensitive data

	return true, nil
}

/*
BindAndValidateURLQueryRequest binds and validates the request in URL query format.
Req should be a pointer to the request struct.
*/
func BindAndValidateURLQueryRequest(c *fiber.Ctx, req interface{}) (bool, error) {
	if err := c.QueryParser(req); err != nil {
		return false, response.Response{
			HttpStatusCode: 400,
			Status:         false,
			ResponseCode:   bindingErrCode,
			Error:          fmt.Errorf("error parsing request:%w", err),
		}.WriteToJSON(c)
	}

	if ok, errResponse := validateRequestInDetailBasedOnFieldTag(c, req, fieldTag_Query); !ok {
		return false, errResponse
	}

	log.Debug("req after validation:", req) //alter later if need to hide sensitive data
	return true, nil
}
