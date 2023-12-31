package api

import (
	"github.com/dalefeng/chat-api-reverse/api/chatgpt"
	"github.com/dalefeng/chat-api-reverse/api/copilot"
)

type ApiGroup struct {
	copilot.CopilotApi
	chatgpt.ChatGPTApi
}

var ApiGroupApp = new(ApiGroup)
