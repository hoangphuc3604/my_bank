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
	UserTableName = "User"
	SessionTableName = "Session"
)

var (
	ErrorSessionBlocked = fmt.Errorf("Session is blocked")
	IncorrectSessionUsername = fmt.Errorf("Incorrect session username")
	ErrorIncorrectRefreshToken = fmt.Errorf("Incorrect refresh token")
	ErrorExpiredToken = fmt.Errorf("Token is expired")
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

func NewCustomError(code int, root error, message string, key string) *AppErr {
	if root != nil {
		return NewFullErrorResponse(code, root, message, root.Error(), key)
	}

	return NewFullErrorResponse(code, errors.New(message), message, message, key)
}

func ErrorDB(err error) *AppErr {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "Something went wrong with DB", err.Error(), "DB_ERROR")
}

func ErrorBinding(err error) *AppErr {
	return NewCustomError(http.StatusInternalServerError, err, "Binding error", "BINDING_ERROR")
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
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("Can not create %s", strings.ToLower(entityName)),
		fmt.Sprintf("CAN_NOT_CREATE_%s", strings.ToUpper(entityName)),
	)
}

func ErrorCanNotGetEntity(entityName string, err error) *AppErr {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("Can not get %s", strings.ToLower(entityName)),
		fmt.Sprintf("CAN_NOT_GET_%s", strings.ToUpper(entityName)),
	)
}

func ErrorCanNotListEntities(entityName string, err error) *AppErr {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("Can not list %s", strings.ToLower(entityName)),
		fmt.Sprintf("CAN_NOT_LIST_%s", strings.ToUpper(entityName)),
	)
}

func ErrorCanNotCountEntities(entityName string, err error) *AppErr {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		fmt.Sprintf("Can not count %s", strings.ToLower(entityName)),
		fmt.Sprintf("CAN_NOT_COUNT_%s", strings.ToUpper(entityName)),
	)
}

func ErrorCanNotTransfer(err error) *AppErr {
	return NewCustomError(http.StatusInternalServerError, err, "Can not transfer", "CAN_NOT_TRANSFER")
}

func ErrorCurrencyMismatch(err error) *AppErr {
	return NewCustomError(
		http.StatusInternalServerError,
		err,
		"Account currency mismatch",
		"ACCOUNT_CURRENCY_MISMATCH",
	)
}

func ErrorDuplicatedEntity(entityName string, err error) *AppErr {
	return NewCustomError(
		http.StatusConflict,
		err,
		fmt.Sprintf("Duplicated %s", strings.ToLower(entityName)),
		fmt.Sprintf("DUPLICATED_%s", strings.ToUpper(entityName)),
	)
}

func ErrorNotExists(entityName string, err error) *AppErr {
	return NewCustomError(
		http.StatusConflict,
		err,
		fmt.Sprintf("%s not exists", strings.ToLower(entityName)),
		fmt.Sprintf("NOT_EXISTS_%s", strings.ToUpper(entityName)),
	)
}

func ErrorAlreadyExists(entityName string, err error) *AppErr {
	return NewCustomError(
		http.StatusConflict,
		err,
		fmt.Sprintf("%s already exists", strings.ToLower(entityName)),
		fmt.Sprintf("ALREADY_EXISTS_%s", strings.ToUpper(entityName)),
	)
}

func ErrorHashPassword(err error) *AppErr {
	return NewCustomError(http.StatusInternalServerError, err, "Can not hash password", "HASH_PASSWORD")
}

func ErrorUnauthorized(err error) *AppErr {
	return NewCustomError(http.StatusUnauthorized, err, "Unauthorized", "UNAUTHORIZED")
}

func ErrorCanNotCreateToken(err error) *AppErr {
	return NewCustomError(http.StatusInternalServerError, err, "Can not create token", "CAN_NOT_CREATE_TOKEN")
}