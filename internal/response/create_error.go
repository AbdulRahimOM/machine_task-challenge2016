package response

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
