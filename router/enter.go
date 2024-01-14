package router

import (
	"github.com/dalefeng/chat-api-reverse/router/chatgpt"
	"github.com/dalefeng/chat-api-reverse/router/copilot"
	"github.com/dalefeng/chat-api-reverse/router/gemini"
)

type Group struct {
	Copilot copilot.Router
	ChatGPT chatgpt.Router
	GeMini  gemini.Router
}

var GroupApp = new(Group)
