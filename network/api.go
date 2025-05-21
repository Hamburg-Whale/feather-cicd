package network

import (
	"feather/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type api struct {
	server *Server
}

func registerServer(server *Server) {
	api := &api{server: server}
	server.engine.POST("/api/v1/user", api.createUser)
}

func (api *api) createUser(ctx *gin.Context) {
	var req types.CreateUserReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if res, err := api.server.service.CreateUser(req.Email, req.Password); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, res)
	}
}
