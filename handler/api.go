package handler

import (
	"feather/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler struct {
	server *Server
}

func registerServer(server *Server) {
	handler := &handler{server: server}
	server.engine.POST("/api/v1/user", handler.createUser)
	server.engine.GET("/api/v1/user/:id", handler.User)

	server.engine.POST("/api/v1/repo", handler.createRepo)

	server.engine.POST("/api/v1/ci", handler.createArgoCi)
}

func (handler *handler) createUser(ctx *gin.Context) {
	var req types.CreateUserReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := handler.server.service.CreateUser(req.Email, req.Password, req.Nickname); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}

func (handler *handler) User(ctx *gin.Context) {
	userId := ctx.Param("id")
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		response(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if res, err := handler.server.service.User(id); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, res)
	}
}

func (handler *handler) createRepo(ctx *gin.Context) {
	var req *types.RepoFromTemplateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if res, err := handler.server.service.CreateRepo(req); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, res)
	}
}

func (handler *handler) createArgoCi(ctx *gin.Context) {
	var req *types.JobBasedJavaRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := handler.server.service.CreateArgoSensor(req); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}
