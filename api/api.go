package api

import (
	"github.com/dalefengs/chat-api-proxy/model/common"
	"github.com/dalefengs/chat-api-proxy/model/common/response"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type Api struct{}

func (a *Api) CompletionsHandler(c *gin.Context) {
	var req openai.CompletionRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("param format error", c)
		return
	}
	switch req.Model {
	case common.CoCopilot:

	}
}
