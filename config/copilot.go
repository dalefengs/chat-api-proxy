package config

type Copilot struct {
	TokenURL       string `yaml:"tokenUrl"`
	CoTokenURL     string `yaml:"coTokenUrl"`
	CompletionsURL string `yaml:"completionsUrl"`
}
