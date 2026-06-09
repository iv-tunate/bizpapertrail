package utils

import (
	"log"

	"github.com/labstack/echo"
)

type ResponseDetail struct {
	Code             int `json:"status_code"`
	Data             any `json:"data"`
	PaginationDetail any `json:"pagination_detail"`
	Message          any `json:"message"`
	Error            any `json:"error"`
}

func ErrorResponse(c echo.Context, code int, message string, err any) error{
	if code > 499{
		log.Printf("Internal Server Error 5xx: %s", message)
	}

	return c.JSON(code, ResponseDetail{
		Code: code,
		Message: message,
		Error: err,
	})
}

func SuccessResponse(c echo.Context, code int, message string, data any, pagination_detail map[string]any) error{
	return c.JSON(code, ResponseDetail{
		Code: code,
		Message: message,
		Data: data,
		PaginationDetail: pagination_detail,
	})
}