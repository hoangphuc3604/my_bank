package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoangphuc3064/MyBank/common"
)

type renewAccessRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessResponse struct {
	AccessToken string `json:"access_token"`
	AccessTokenExpires time.Time `json:"access_token_expires"`
}

func (server *Server) renewAccess(ctx *gin.Context) {
	var request renewAccessRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(common.ErrorBinding(err)))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse(common.ErrorUnauthorized(err)))
		return
	}

	session, err := server.store.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse(common.ErrorNotFound(common.SessionTableName)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotGetEntity(common.SessionTableName, err)))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse(common.ErrorUnauthorized(common.ErrorSessionBlocked)))
		return
	}

	if session.Username != payload.Username {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse(common.ErrorUnauthorized(common.IncorrectSessionUsername)))
		return
	}

	if session.RefreshToken != request.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse(common.ErrorUnauthorized(common.ErrorIncorrectRefreshToken)))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse(common.ErrorUnauthorized(common.ErrorExpiredToken)))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(payload.Username, "User", server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotCreateToken(err)))
		return
	}

	response := renewAccessResponse{
		AccessToken: accessToken,
		AccessTokenExpires: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(response))
}