package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Err struct {
	Err    error
	Status int
	Msg    string
}

func (e *Err) Error() string {
	return fmt.Sprintf("%s (status: %d)", e.Msg, e.Status)
}

func NewErr(err error, status int, msg string) error {
	return &Err{
		Err:    err,
		Status: status,
		Msg:    msg,
	}
}

func WriteErr(c *gin.Context, err error) {
	var syntaxErr *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var validationErrors validator.ValidationErrors
	var appErr *Err

	if errors.As(err, &appErr) {
		if appErr.Err != nil {
			logrus.Printf("Internal error: %v", appErr.Err)
		}
		newErrorResponse(c, appErr.Status, map[string]string{"Error": appErr.Msg})
		return
	}

	if errors.As(err, &syntaxErr) {
		msg := fmt.Sprintf("Syntax error Json. Position: %v", syntaxErr.Offset)
		newErrorResponse(c, http.StatusBadRequest, map[string]string{"Error": msg})
		return
	}

	if errors.As(err, &unmarshalTypeError) {
		msg := fmt.Sprintf("Invalid field data type: %s . Expected: %s", unmarshalTypeError.Field, unmarshalTypeError.Type)
		newErrorResponse(c, http.StatusBadRequest, map[string]string{"Error": msg})
	}

	if errors.Is(err, io.EOF) {
		newErrorResponse(c, http.StatusBadRequest, map[string]string{"Error": "Empty request body"})
		return
	}

	if errors.As(err, &validationErrors) {
		newErrorResponse(c, http.StatusBadRequest, map[string]string{"Error": "Data validation error"})
		return
	}

	logrus.Printf("Unknown error: %v", err)
	newErrorResponse(c, http.StatusInternalServerError, map[string]string{"Error": "Something happened.."})

}

func newErrorResponse(c *gin.Context, statusCode int, data interface{}) {
	logrus.Error(data)
	c.AbortWithStatusJSON(statusCode, data)
}
