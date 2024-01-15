package router

import (
	"github.com/dalefengs/chat-api-proxy/router/copilot"
	"github.com/dalefengs/chat-api-proxy/router/genai"
	"github.com/dalefengs/chat-api-proxy/router/openai"
	"github.com/dalefengs/chat-api-proxy/router/token"
)

type Group struct {
	Copilot copilot.Router
	ChatGPT openai.Router
	GeMini  genai.Router
	Token   token.Router
}

var GroupApp = new(Group)
