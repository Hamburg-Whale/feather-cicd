package handler

import (
	"feather/service"
	"feather/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BasecampHandler struct {
	basecampService service.BasecampService
}

func NewBasecampHandler(bs service.BasecampService) *BasecampHandler {
	return &BasecampHandler{basecampService: bs}
}

func (h *BasecampHandler) CreateBasecamp(ctx *gin.Context) {
	var req types.CreateBasecampReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := h.basecampService.CreateBasecamp(req.Name, req.URL, req.Token, req.Owner, req.User_ID); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}

func (h *BasecampHandler) Basecamp(ctx *gin.Context) {
	basecampId := ctx.Param("id")
	id, err := strconv.ParseInt(basecampId, 10, 64)
	if err != nil {
		response(ctx, http.StatusBadRequest, "Invalid basecamp ID")
		return
	}

	if res, err := h.basecampService.Basecamp(id); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, res)
	}
}

func (h *BasecampHandler) BasecampByUserId(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		response(ctx, http.StatusBadRequest, "Invalid basecamp ID")
		return
	}

	if res, err := h.basecampService.BasecampsByUserId(id); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, res)
	}
}
