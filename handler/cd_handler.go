package handler

import (
	"feather/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CdHandler struct {
	argoCdService service.ArgoCdService
}

func NewCdHandler(as service.ArgoCdService) *CdHandler {
	return &CdHandler{argoCdService: as}
}

func (h *CdHandler) CreateArgoCd(ctx *gin.Context) {
	baseCampId := ctx.Param("id")

	id, err := strconv.ParseInt(baseCampId, 10, 64)
	if err != nil {
		response(ctx, http.StatusBadRequest, "Invalid BaseCamp ID")
		return
	}

	if err := h.argoCdService.CreateProjectManifestRepo(id); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}
