package router

import (
	"github.com/dalefengs/chat-api-proxy/router/copilot"
	"github.com/dalefengs/chat-api-proxy/router/genai"
	"github.com/dalefengs/chat-api-proxy/router/openai"
)

type Group struct {
	Copilot copilot.Router
	ChatGPT openai.Router
	GeMini  genai.Router
}

var GroupApp = new(Group)
