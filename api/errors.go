package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
		if apiError, ok := err.(Error); ok {
			return c.Status(apiError.Code).JSON(apiError)
		}
		apiError := NewError(http.StatusInternalServerError, err.Error())
		return c.Status(apiError.Code).JSON(apiError)
	}

//Error implements the Error interface
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrUnauthortized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err: "unauthorized request",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "invalid id given",
	}
}

func ErrBadRequest() Error{
	return Error{
		Code:http.StatusBadRequest,
		Err: "Invalid JSON request",
	}
}

func ErrResourceNotFound(res string) Error{
	return Error{
		Code:http.StatusNotFound,
		Err: res + " Resource not found",
	}
}