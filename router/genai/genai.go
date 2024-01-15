package genai

import (
	"github.com/dalefengs/chat-api-proxy/api/genai"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (s *Router) InitGeminiRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("gemini")
	genApi := genai.GenApi{}
	{
		baseRouter.POST("/v1/chat/completions", genApi.CompletionsHandler)
	}
	return baseRouter
}
