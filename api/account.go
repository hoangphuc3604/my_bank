package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangphuc3064/MyBank/common"
	"github.com/hoangphuc3064/MyBank/db/sqlc"
	"github.com/hoangphuc3064/MyBank/token"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var request createAccountRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(common.ErrorBinding(err)))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := sqlc.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: request.Currency,
		Balance:  0,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, common.ErrorResponse(common.ErrorNotExists(common.UserTableName, err)))
				return

			case "unique_violation":
				ctx.JSON(http.StatusConflict, common.ErrorResponse(common.ErrorAlreadyExists(common.AccountTableName, err)))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotCreateEntity(common.AccountTableName, err)))
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
}

func (server *Server) getAccount(ctx *gin.Context) {
	var request getAccountRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(common.ErrorBinding(err)))
	}

	acc, err := server.store.GetAccount(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse(common.ErrorNotFound(common.AccountTableName)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotGetEntity(common.AccountTableName, err)))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if acc.Owner != authPayload.Username {
		ctx.JSON(http.StatusForbidden, common.ErrorResponse(common.ErrorUnauthorized(fmt.Errorf("account does not belong to the authenticated user"))))
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(acc))
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var request common.Paging
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(common.ErrorBinding(err)))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := sqlc.ListAccountsParams{
		Owner: authPayload.Username,
		Limit:  int32(request.Limit),
		Offset: int32(request.Offset()),
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotListEntities(common.AccountTableName, err)))
		return
	}
	request.Total, err = server.store.CountAccounts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotCountEntities(common.AccountTableName, err)))
		return
	}

	ctx.JSON(http.StatusOK, common.NewSuccessResponse(accounts, request, nil))
}