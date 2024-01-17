package api

import (
	"github.com/dalefengs/chat-api-proxy/api"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (s *Router) InitTokenRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("")
	apiApi := api.Api{}
	{
		baseRouter.POST("/v1/chat/completions", apiApi.CompletionsHandler)
	}
	return baseRouter
}
