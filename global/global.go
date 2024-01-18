package global

import (
	"github.com/dalefengs/chat-api-proxy/global/config"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Config           config.Server
	Viper            *viper.Viper
	Log              *zap.Logger
	SugarLog         *zap.SugaredLogger
	UserHomeDir      string // 家目录
	UserHomeCacheDir string // 家缓存目录
)

var Json = jsoniter.ConfigCompatibleWithStandardLibrary
