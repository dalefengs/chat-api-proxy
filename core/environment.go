package core

import (
	"github.com/dalefengs/chat-api-proxy/global"
	"os"
)

// InitEnvironmentVar 初始化环境参数
func InitEnvironmentVar() {
	if env := os.Getenv("PROXY_API_PREFIX"); env != "" {
		global.Config.System.RouterPrefix = env
	}
	if env := os.Getenv("LOG_LEVEL"); env != "" {
		global.Config.Zap.Level = env
	}
	if env := os.Getenv("GEMINI_BASE_URL"); env != "" {
		global.Config.Gemini.BaseUrl = env
	}
	if env := os.Getenv("GEMINI_VERSION"); env != "" {
		global.Config.Gemini.ApiVersion = env
	}
}
