package api

import (
	"github.com/dalefeng/chat-api-reverse/api/chatgpt"
	"github.com/dalefeng/chat-api-reverse/api/copilot"
	"github.com/dalefeng/chat-api-reverse/api/gemini"
)

type ApiGroup struct {
	copilot.CopilotApi
	chatgpt.ChatGPTApi
	gemini.GeMiniApi
}

var ApiGroupApp = new(ApiGroup)
