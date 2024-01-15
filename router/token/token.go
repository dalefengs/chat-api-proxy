package token

import (
	"github.com/dalefengs/chat-api-proxy/api/token"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (s *Router) InitTokenRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("token")
	tokenApi := token.TokenApi{}
	{
		baseRouter.POST("/pool", tokenApi.TokenPoolHandler)
	}
	return baseRouter
}
