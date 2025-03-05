package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangphuc3064/MyBank/common"
	"github.com/hoangphuc3064/MyBank/db/sqlc"
	"github.com/hoangphuc3064/MyBank/token"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var request createTransferRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(common.ErrorBinding(err)))
		return
	}

	fromAcc, valid := server.validateAccount(ctx, request.FromAccountID, request.Currency)
	if !valid {
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAcc.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse(common.ErrorUnauthorized(fmt.Errorf("account [%d] does not belong to user [%s]", request.FromAccountID, authPayload.Username))))
		return
	}

	_, valid = server.validateAccount(ctx, request.ToAccountID, request.Currency)
	if !valid {
		return
	}

	arg := sqlc.TransferTxParams{
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		Amount:        request.Amount,
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotTransfer(err)))
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
}

func (server *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) (sqlc.Account,bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse(common.ErrorNotFound(common.AccountTableName)))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotGetEntity(common.AccountTableName, err)))
		return account, false
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(common.ErrorCurrencyMismatch(
			fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency),
		)))
		return account, false
	}

	return account, true
}