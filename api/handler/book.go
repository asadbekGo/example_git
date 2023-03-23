package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateBook(c *gin.Context) {

	var createBook models.CreateBook

	err := c.ShouldBindJSON(&createBook)
	if err != nil {
		h.handlerResponse(c, "create book", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storages.Book().Create(&createBook)
	if err != nil {
		h.handlerResponse(c, "storage.book.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.storages.Book().GetByID(&models.BookPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.book.getByID", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create book", http.StatusCreated, resp)
}

func (h *Handler) GetByIdBook(c *gin.Context) {

	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "get by id book", http.StatusBadRequest, "invalid book id")
		return
	}

	resp, err := h.storages.Book().GetByID(&models.BookPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.book.getByID", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create book", http.StatusCreated, resp)
}

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
