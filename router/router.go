package router

import (
	"github.com/dalefengs/chat-api-proxy/api"
	"github.com/gin-gonic/gin"
)

func (r *Router) InitRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("")
	copilotApi := api.ApiGroupApp.CopilotApi
	_ = api.ApiGroupApp.CommonApi
	{
		baseRouter.GET("/copilot_internal/v2/token", copilotApi.CoTokenHandler)
		baseRouter.POST("/v1/chat/completions", copilotApi.CompletionsHandler)
		baseRouter.POST("/chat/completions", copilotApi.CompletionsOfficialHandler)
	}
	return baseRouter
}
