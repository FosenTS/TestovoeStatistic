package http

import (
	"Testovoe1/internal/storage"
	ginresponse "Testovoe1/pkg/ginresponce"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"net/http"
	"strconv"
	"time"
)

type HttpHandlerController struct {
	ginEngine *gin.Engine
	storage   *storage.Storage
	logger    log.Logger
}

func NewHttpHandlerController(storage *storage.Storage, logger log.Logger) *HttpHandlerController {
	return &HttpHandlerController{
		ginEngine: gin.Default(),
		storage:   storage,
		logger:    logger,
	}
}

func (h *HttpHandlerController) RegisterRouter() {
	h.ginEngine.GET("/posts/comments/statistics", h.GetStatistics)
}

func (h *HttpHandlerController) StartHTTP() {
	h.RegisterRouter()
	h.ginEngine.Run()
}

func (h *HttpHandlerController) GetStatistics(c *gin.Context) {
	logger := log.With(h.logger, "method", "GetStatistics")
	postIdS, ok := c.GetQuery("postId")

	if ok {
		postId, err := strconv.Atoi(postIdS)
		if err != nil {
			level.Error(logger).Log("err", err)
			ginresponse.ErrorString(c, http.StatusInternalServerError, err, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(c, time.Second*40)
		defer cancel()
		resEntity, err := h.storage.GetStatististicById(ctx, postId)
		if err != nil {
			level.Error(logger).Log("err", err)
			ginresponse.ErrorString(c, http.StatusInternalServerError, err, err.Error())
			return
		}
		c.JSON(http.StatusOK, resEntity)
	} else {
		ctx, cancel := context.WithTimeout(c, time.Second*40)
		defer cancel()
		resEntity, err := h.storage.GetAllStatistics(ctx)
		if err != nil {
			level.Error(logger).Log("err", err)
			ginresponse.ErrorString(c, http.StatusInternalServerError, err, err.Error())
			return
		}
		c.JSON(http.StatusOK, resEntity)
	}

}
