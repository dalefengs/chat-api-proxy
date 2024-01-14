package gemini

import (
	"github.com/dalefeng/chat-api-reverse/api/gemini"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (s *Router) InitGeminiRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("gemini")
	geminiApi := gemini.GeMiniApi{}
	{
		baseRouter.POST("/v1/chat/completions", geminiApi.CompletionsHandler)
	}
	return baseRouter
}
