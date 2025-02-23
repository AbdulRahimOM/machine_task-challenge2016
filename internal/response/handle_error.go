package response

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
)

var (
	sqlRegexPattern = regexp.MustCompile(`SQLSTATE (\d{5})`)
)

type custError struct {
	Response
	Error string `json:"error"`
}

func (resp Response) WriteToJSON(c *fiber.Ctx) error {
	if resp.Error == nil {
		return c.Status(resp.HttpStatusCode).JSON(resp)
	}
	newCustError := custError{
		Response: resp,
	}
	if resp.Error != nil {
		newCustError.Error = resp.Error.Error()
	}

	return c.Status(resp.HttpStatusCode).JSON(newCustError)
}
