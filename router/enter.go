package router

import (
	"github.com/dalefeng/chat-api-reverse/router/chatgpt"
	"github.com/dalefeng/chat-api-reverse/router/copilot"
)

type Group struct {
	Copilot copilot.Router
	ChatGPT chatgpt.Router
}

var GroupApp = new(Group)
