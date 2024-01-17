package common

const (
	ChatGPT   = "chatgpt"
	Gemini    = "gemini"
	Copilot   = "copilot"
	CoCopilot = "cocopilot"
)

var Models = map[string]map[string]string{
	Gemini: {
		"gemini-pro":        "gemini-pro",
		"gemini-pro-vision": "gemini-pro-vision",
	},
	Copilot: {
		"gpt-4": "gpt4",
	},
}
