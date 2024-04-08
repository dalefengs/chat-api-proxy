package config

type System struct {
	Env          string `mapstructure:"env" json:"env" yaml:"env"`    // 环境值
	Port         int    `mapstructure:"port" json:"port" yaml:"port"` // 端口值
	RouterPrefix string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
}

type Proxy struct {
	HTTP  string `mapstructure:"http" json:"http" yaml:"http"`
	HTTPS string `mapstructure:"https" json:"https" yaml:"https"`
}
