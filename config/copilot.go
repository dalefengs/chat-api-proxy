package config

type Copilot struct {
	TokenURL       string `yaml:"tokenUrl"`
	CopilotHost    string `yaml:"copilotHost"`
	CoCopilotHost  string `yaml:"coCopilotHost"`
	CoTokenURL     string `yaml:"coTokenUrl"`
	CompletionsURL string `yaml:"completionsUrl"`
}
