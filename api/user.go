package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hoangphuc3064/MyBank/common"
	"github.com/hoangphuc3064/MyBank/db/sqlc"
	"github.com/hoangphuc3064/MyBank/util"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password   string `json:"password" binding:"required,min=6"`
	Fullname   string `json:"fullname" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username  string    `json:"username"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserResponse(user sqlc.User) *userResponse {
	return &userResponse{
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var request createUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(common.ErrorBinding(err)))
		return
	}

	hashPassword, err := util.HashPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorHashPassword(err)))
		return
	}

	arg := sqlc.CreateUserParams{
		Username: request.Username,
		Password: hashPassword,
		Fullname: request.Fullname,
		Email: request.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, common.ErrorResponse(common.ErrorAlreadyExists(common.UserTableName, err)))
				return
			default:
				ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotCreateEntity(common.UserTableName, err)))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotCreateEntity(common.UserTableName, err)))
		return
	}

	response := NewUserResponse(user)
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(response))
}

type loginRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password   string `json:"password" binding:"required,min=6"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	AccessTokenExpires time.Time `json:"access_token_expires"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpires time.Time `json:"refresh_token_expires"`
	User 	  userResponse `json:"user"`
	SessionID uuid.UUID `json:"session_id"`
}

func (server *Server) login(ctx *gin.Context) {
	var request loginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(common.ErrorBinding(err)))
		return
	}

	user, err := server.store.GetUser(ctx, request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, common.ErrorResponse(common.ErrorNotFound(common.UserTableName)))
			return
		}

		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotGetEntity(common.UserTableName, err)))
		return
	}

	err = util.CheckPassword(request.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, common.ErrorResponse(common.ErrorUnauthorized(err)))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, "User", server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotCreateToken(err)))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, "User", server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotCreateToken(err)))
		return
	}

	session, err := server.store.CreateSession(ctx, sqlc.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:    user.Username,
		RefreshToken: refreshToken,
		UserAgent:   ctx.Request.UserAgent(),
		IpAddress:  ctx.ClientIP(),
		IsBlocked: false,
		ExpiresAt: refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotCreateEntity(common.SessionTableName, err)))
		return
	}

	response := loginResponse{
		AccessToken: accessToken,
		AccessTokenExpires: accessPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpires: refreshPayload.ExpiredAt,
		User: *NewUserResponse(user),
		SessionID: session.ID,
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(response))
}