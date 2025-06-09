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

	server.engine.POST("/api/v1/basecamp", handler.createBasecamp)
	server.engine.GET("/api/v1/basecamp/:id", handler.BaseCamp)

	server.engine.POST("/api/v1/project", handler.CreateProject)
	server.engine.GET("/api/v1/project/:id", handler.Project)

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

func (handler *handler) createBasecamp(ctx *gin.Context) {
	var req types.CreateBasecampReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := handler.server.service.CreateBaseCamp(req.Name, req.URL, req.Token, req.User_ID); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}

func (handler *handler) BaseCamp(ctx *gin.Context) {
	basecampId := ctx.Param("id")
	id, err := strconv.ParseInt(basecampId, 10, 64)
	if err != nil {
		response(ctx, http.StatusBadRequest, "Invalid basecamp ID")
		return
	}

	if res, err := handler.server.service.BaseCamp(id); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, res)
	}
}

func (handler *handler) CreateProject(ctx *gin.Context) {
	var req types.CreateProjectReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := handler.server.service.CreateProject(req.Name, req.URL, req.Owner, req.Private, req.BaseCamp_ID); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}

func (handler *handler) Project(ctx *gin.Context) {
	projectId := ctx.Param("id")
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		response(ctx, http.StatusBadRequest, "Invalid project ID")
		return
	}

	if res, err := handler.server.service.Project(id); err != nil {
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
