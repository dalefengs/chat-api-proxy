package chatgpt

import (
	"github.com/dalefeng/chat-api-reverse/api"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (s *Router) InitChatGPTRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("")
	chaGPTApi := api.ApiGroupApp.ChatGPTApi
	{
		baseRouter.POST("v1/chat/completions", chaGPTApi.Completions)
	}
	return baseRouter
}
