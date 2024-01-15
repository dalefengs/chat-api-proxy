package openai

import (
	"github.com/dalefeng/chat-api-reverse/api"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (s *Router) InitChatGPTRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("")
	openApi := api.ApiGroupApp.OpenApi
	{
		baseRouter.POST("v1/chat/completions", openApi.Completions)
	}
	return baseRouter
}
