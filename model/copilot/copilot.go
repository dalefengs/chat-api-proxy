package copilot

type Config struct {
	TokenUrl       string `yaml:"tokenUrl"`
	CoTokenUrl     string `yaml:"cotokenUrl"`
	CompletionsUrl string `yaml:"completionsUrl"`
}
