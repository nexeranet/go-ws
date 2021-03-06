package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nexeranet/go-ws/pkg/ws"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (*Handler) InitRouter(pWS *ws.WS, hub *ws.Hub) *gin.Engine {
	router := gin.New()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.tmpl")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
	router.GET("/ws", func(c *gin.Context) {
		ws.WsHandler(c.Writer, c.Request, pWS, hub)
	})
	return router
}
