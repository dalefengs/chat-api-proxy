package api

import (
	"github.com/dalefengs/chat-api-proxy/api/copilot"
	"github.com/dalefengs/chat-api-proxy/api/genai"
	"github.com/dalefengs/chat-api-proxy/api/openai"
	"github.com/dalefengs/chat-api-proxy/api/token"
)

type ApiGroup struct {
	token.TokenApi
	copilot.CopilotApi
	openai.OpenApi
	genai.GenApi
	Api
}

var ApiGroupApp = new(ApiGroup)
