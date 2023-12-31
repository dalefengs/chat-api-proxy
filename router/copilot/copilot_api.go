package copilot

import (
	"github.com/dalefeng/chat-api-reverse/api"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (s *Router) InitCopilotRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("copilot")
	copilotApi := api.ApiGroupApp.CopilotApi
	{
		baseRouter.POST("v1/chat/completions", copilotApi.Completions)
	}
	return baseRouter
}
