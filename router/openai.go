package router

import (
	"github.com/dalefengs/chat-api-proxy/api"
	"github.com/gin-gonic/gin"
)

func (r *Router) InitChatGPTRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("")
	openApi := api.ApiGroupApp.OpenApi
	{
		baseRouter.POST("v1/chat/completions", openApi.CompletionsHandler)
	}
	return baseRouter
}
