package api

import (
	"github.com/dalefeng/chat-api-reverse/api/copilot"
	"github.com/dalefeng/chat-api-reverse/api/genai"
	"github.com/dalefeng/chat-api-reverse/api/openai"
)

type ApiGroup struct {
	copilot.CopilotApi
	openai.OpenApi
	genai.GenApi
}

var ApiGroupApp = new(ApiGroup)
