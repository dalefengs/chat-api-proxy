package config

type System struct {
	Env          string `mapstructure:"env" json:"env" yaml:"env"`    // 环境值
	Port         int    `mapstructure:"port" json:"port" yaml:"port"` // 端口值
	RouterPrefix string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
}
