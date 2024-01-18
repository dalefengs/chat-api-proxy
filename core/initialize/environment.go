package initialize

import (
	"github.com/dalefengs/chat-api-proxy/global"
	"os"
)

func InitEnvironmentVar() {
	if env := os.Getenv("GEMINI_BASE_URL"); env != "" {
		global.Config.Gemini.BaseUrl = baseUrl
	}
	if env := os.Getenv("GEMINI_VERSION"); env != "" {
		global.Config.Gemini.BaseUrl = env
	}
}
