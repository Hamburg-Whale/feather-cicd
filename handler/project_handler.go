package handler

import (
	"feather/service"
	"feather/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(ps service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: ps}
}

func (h *ProjectHandler) CreateProject(ctx *gin.Context) {
	var req types.CreateProjectReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := h.projectService.CreateProject(req.Name, req.URL, req.Owner, req.Private, req.BaseCamp_ID); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}

func (h *ProjectHandler) Project(ctx *gin.Context) {
	projectId := ctx.Param("id")
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		response(ctx, http.StatusBadRequest, "Invalid project ID")
		return
	}

	if res, err := h.projectService.Project(id); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, res)
	}
}
