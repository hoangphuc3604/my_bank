package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
// type getUserRequest struct {
// 	Username string `uri:"username" binding:"required,alphanum"`
// }

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
			}
		}

		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotCreateEntity(common.UserTableName, err)))
		return
	}

	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(user))
}

// func (server *Server) getUser(ctx *gin.Context) {
// 	var request getUserRequest
// 	if err := ctx.ShouldBindUri(&request); err != nil {
// 		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(common.ErrorBinding(err)))
// 	}

// 	acc, err := server.store.GetAccount(ctx, request.ID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, common.ErrorResponse(common.ErrorNotFound(common.AccountTableName)))
// 			return
// 		}

// 		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotGetEntity(common.AccountTableName, err)))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(acc))
// }

// func (server *Server) listUsers(ctx *gin.Context) {
// 	var request common.Paging
// 	if err := ctx.ShouldBindQuery(&request); err != nil {
// 		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(common.ErrorBinding(err)))
// 		return
// 	}

// 	arg := sqlc.ListAccountsParams{
// 		Limit:  int32(request.Limit),
// 		Offset: int32(request.Offset()),
// 	}

// 	accounts, err := server.store.ListAccounts(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotListEntities(common.AccountTableName, err)))
// 		return
// 	}
// 	request.Total, err = server.store.CountAccounts(ctx)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(common.ErrorCanNotCountEntities(common.AccountTableName, err)))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, common.NewSuccessResponse(accounts, request, nil))
// }