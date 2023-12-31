package copilot

import (
	"github.com/dalefeng/chat-api-reverse/model/common/response"
	"github.com/gin-gonic/gin"
)

type CopilotApi struct {
}

func (co *CopilotApi) Completions(c *gin.Context) {
	response.OkWithMessage("success", c)
}
