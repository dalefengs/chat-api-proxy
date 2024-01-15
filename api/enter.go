package api

import (
	"github.com/dalefengs/chat-api-proxy/api/copilot"
	"github.com/dalefengs/chat-api-proxy/api/genai"
	"github.com/dalefengs/chat-api-proxy/api/openai"
)

type ApiGroup struct {
	copilot.CopilotApi
	openai.OpenApi
	genai.GenApi
}

var ApiGroupApp = new(ApiGroup)
