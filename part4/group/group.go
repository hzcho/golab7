package group

import (
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	log *logrus.Logger
}

func NewHandler(log *logrus.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	api := router.Group("/api/v1", h.mid)
	{
		api.GET("/hello", h.Get)
		api.POST("/data", h.Post)
		api.PUT("/list", h.Put)
		api.DELETE("/list", h.Delete)
	}
}

func (h *Handler) mid(c *gin.Context) {
	start := time.Now()

	c.Next()

	duration := time.Since(start)
	h.log.WithFields(logrus.Fields{
		"url":      c.Request.URL.String(),
		"method":   c.Request.Method,
		"duration": duration,
	}).Info("Запрос обработан")
}

func (h *Handler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, "hello world")
}

func (h *Handler) Post(c *gin.Context) {
	var json map[string]interface{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info(json)
}

func (h *Handler) Put(c *gin.Context) {
	h.log.Info("list is updated")
}

func (h *Handler) Delete(c *gin.Context) {
	h.log.Info("list is deleted")
}
