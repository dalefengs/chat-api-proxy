package router

import (
	"github.com/dalefeng/chat-api-reverse/router/copilot"
	"github.com/dalefeng/chat-api-reverse/router/genai"
	"github.com/dalefeng/chat-api-reverse/router/openai"
)

type Group struct {
	Copilot copilot.Router
	ChatGPT openai.Router
	GeMini  genai.Router
}

var GroupApp = new(Group)
