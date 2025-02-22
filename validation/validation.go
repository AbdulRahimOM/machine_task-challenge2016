package validation

import (
	"challenge16/internal/response"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gofiber/fiber/v2/log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func validateRequestInDetailBasedOnFieldTag(c *fiber.Ctx, req interface{}, tag string) (bool, error) {

	errorResponses := []response.InvalidField{}
	errs := validate.Struct(req)

	if errs != nil {

		for _, err := range errs.(validator.ValidationErrors) {
			// Get the required tag name using reflection
			formFieldKey := getFieldKeyByStructTag(req, err.Field(), tag)

			e := response.InvalidField{
				FailedField: formFieldKey,
				Tag:         err.Tag(),
				Value:       err.Value(),
			}

			// switch e.FailedField {
			// case "password", "Password":
			// 	log.Debug(fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'", e.FailedField, "--hidden--", e.Tag))
			// default:
			log.Debug(fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'", e.FailedField, e.Value, e.Tag))
			// }

			errorResponses = append(errorResponses, e)
		}
		log.Debug("error validating request:", errorResponses)
		return false, c.Status(http.StatusBadRequest).JSON(response.ValidationErrorResponse{
			Status:       false,
			ResponseCode: validationErrCode,
			Errors:       errorResponses,
		})
	}

	return true, nil
}

// Function to get the tag-name of a struct field based on the tag passed
func getFieldKeyByStructTag(req interface{}, fieldName string, tag string) string {
	val := reflect.TypeOf(req)

	// Check if the value passed is a pointer and get the element type
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Find the struct field by name and return the field's tag based key
	field, found := val.FieldByName(fieldName)
	if !found {
		return fieldName // Return the field name if no such tag is found
	}

	fieldKey := field.Tag.Get(tag)
	if fieldKey == "" {
		return fieldName // Return the field name itself if the given tag is not defined for the field
	}
	return fieldKey
}
