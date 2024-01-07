package copilot

import (
	"github.com/dalefeng/chat-api-reverse/utils"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func GetCompletionsHeader(token string) map[string]string {
	uid := uuid.New().String()
	headersMap := map[string]string{
		"Host":                   "api.githubcopilot.com",
		"Accept-Encoding":        "gzip, deflate, br",
		"Accept":                 "*/*",
		"Authorization":          "Bearer " + token,
		"X-Request-Id":           uid,
		"X-Github-Api-Version":   "2023-07-07",
		"Vscode-Sessionid":       uid + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10),
		"vscode-machineid":       utils.GenHexStr(64),
		"Editor-Version":         "vscode/1.85.0",
		"Editor-Plugin-Version":  "copilot-chat/0.11.1",
		"Openai-Organization":    "github-copilot",
		"Copilot-Integration-Id": "vscode-chat",
		"Openai-Intent":          "conversation-panel",
		"User-Agent":             "GitHubCopilotChat/0.11.1",
	}
	return headersMap
}
