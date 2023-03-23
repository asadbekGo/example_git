package handler

import (
	"app/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetListBook(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list book", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list book", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.storages.Book().GetList(&models.GetListBookRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.book.getlist", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list book response", http.StatusOK, resp)
}
