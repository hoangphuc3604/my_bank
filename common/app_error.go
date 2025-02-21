package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AccountTableName = "Account"
	TransferTableName = "Transfer"
	EntryTableName = "Entry"
)

type AppErr struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"key"`
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err}
}

func NewFullErrorResponse(statusCode int, rootErr error, message string, log string, key string) *AppErr {
	return &AppErr{
		StatusCode: statusCode,
		RootErr:    rootErr,
		Message:    message,
		Log:        log,
		Key:        key,
	}
}

func BadRequestResponse(rootErr error, message string, log string, key string) *AppErr {
	return NewFullErrorResponse(http.StatusBadRequest, rootErr, message, log, key)
}

func (e *AppErr) RootError() error {
	if err, ok := e.RootErr.(*AppErr); ok {
		return err.RootError()
	}

	return e.RootErr
}
func (e *AppErr) Error() string {
	return e.RootError().Error()
}

func NewCustomError(root error, message string, key string) *AppErr {
	if root != nil {
		return BadRequestResponse(root, message, root.Error(), key)
	}

	return BadRequestResponse(errors.New(message), message, message, key)
}

func ErrorDB(err error) *AppErr {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "Something went wrong with DB", err.Error(), "DB_ERROR")
}

func ErrorBinding(err error) *AppErr {
	return NewCustomError(err, "Binding error", "BINDING_ERROR")
}

func ErrorNotFound(entityName string) *AppErr {
	return NewFullErrorResponse(
		http.StatusNotFound,
		fmt.Errorf("%s not found", entityName),
		fmt.Sprintf("%s not found", entityName),
		fmt.Sprintf("%s not found", entityName),
		fmt.Sprintf("NOT_FOUND_%s", strings.ToUpper(entityName)),
	)
}

func ErrorCanNotCreateEntity(entityName string, err error) *AppErr {
	return NewCustomError(
		err,
		fmt.Sprintf("Can not create %s", strings.ToLower(entityName)),
		fmt.Sprintf("CAN_NOT_CREATE_%s", strings.ToUpper(entityName)),
	)
}

func ErrorCanNotGetEntity(entityName string, err error) *AppErr {
	return NewCustomError(
		err,
		fmt.Sprintf("Can not get %s", strings.ToLower(entityName)),
		fmt.Sprintf("CAN_NOT_GET_%s", strings.ToUpper(entityName)),
	)
}

func ErrorCanNotListEntities(entityName string, err error) *AppErr {
	return NewCustomError(
		err,
		fmt.Sprintf("Can not list %s", strings.ToLower(entityName)),
		fmt.Sprintf("CAN_NOT_LIST_%s", strings.ToUpper(entityName)),
	)
}

func ErrorCanNotCountEntities(entityName string, err error) *AppErr {
	return NewCustomError(
		err,
		fmt.Sprintf("Can not count %s", strings.ToLower(entityName)),
		fmt.Sprintf("CAN_NOT_COUNT_%s", strings.ToUpper(entityName)),
	)
}