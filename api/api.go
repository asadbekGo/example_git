package api

import (
	"app/api/handler"
	"app/config"
	"app/pkg/logger"
	"app/storage"

	"github.com/gin-gonic/gin"
)

func NewApi(r *gin.Engine, cfg *config.Config, store storage.StorageI, logger logger.LoggerI) {

	handler := handler.NewHandler(cfg, store, logger)

	r.GET("/book", handler.GetListBook)

}
