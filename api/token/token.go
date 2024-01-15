package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/dalefengs/chat-api-proxy/global"
	"github.com/dalefengs/chat-api-proxy/model/common"
	"github.com/dalefengs/chat-api-proxy/model/common/response"
	tokenModel "github.com/dalefengs/chat-api-proxy/model/token"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type TokenApi struct{}

var TokenCache *bigcache.BigCache

func init() {
	var err error
	config := bigcache.Config{
		// 分片数量（必须是2的幂）
		Shards: 1024,
		// 条目可以被驱逐的时间 7天
		LifeWindow: 30 * 24 * time.Hour,
		// 清除过期条目的时间间隔（清理）。
		// 如果设置为 <= 0，则不执行任何操作。
		// 设置为 < 1 秒会适得其反——bigcache 的分辨率为一秒。
		CleanWindow: 5 * time.Minute,
		// rps * lifeWindow，仅用于初始内存分配
		MaxEntriesInWindow: 1000 * 10 * 60,
		// 最大条目大小（以字节为单位），仅用于初始内存分配
		MaxEntrySize: 500,
		// 有关其他内存分配的信息
		Verbose: true,
		// 缓存不会分配超过此限制的内存，单位为 MB
		// 如果达到该值，则可以覆盖最旧的条目以获取新条目
		// 0 值表示没有大小限制
		HardMaxCacheSize: 0,
		// 当最旧的条目因其过期时间或没有空间留给新条目而被删除时触发的回调，或者是因为调用了 delete。
		// 将返回一个表示原因的位掩码。
		// 默认值为 nil，表示没有回调，并且它可以防止解包最旧的条目。
		OnRemove: nil,
		// OnRemoveWithReason 是当最旧的条目因其过期时间或没有空间留给新条目而被删除时触发的回调，或者是因为调用了 delete。
		// 将传递一个表示原因的常量。
		// 默认值为 nil，表示没有回调，并且它可以防止解包最旧的条目。
		// 如果指定了 OnRemove，则忽略。
		OnRemoveWithReason: nil,
	}

	TokenCache, err = bigcache.New(context.Background(), config)
	if err != nil {
		log.Println("init TokenCache error ", err.Error())
		panic(err)
	}
	log.Println("init TokenCache success")
}

func (a *TokenApi) TokenPoolHandler(c *gin.Context) {
	var req *tokenModel.TokensRequest
	err := c.ShouldBindJSON(req)
	if err != nil {
		response.FailWithMessage("json bind error", c)
		return
	}
	if msg := CheckTokenParam(req); msg != "" {
		response.FailWithMessage(msg, c)
		return
	}
	var token string
	if req.Token != "" {
		token = req.Token
	} else {
		token = GenerateToken("pk")
	}

	tokenCache, err := TokenCache.Get(token)
	if errors.Is(err, bigcache.ErrEntryNotFound) {
		// 设置新 token
		setErr := SetNewTokenInfo(token, req)
		if setErr != nil {
			response.FailWithMessage("failed to save new token pool", c)
			return
		}
	} else if err != nil {
		response.FailWithMessage("failed to find token", c)
		return
	}
	// 更新 poolToken
	println(tokenCache)

}

// SetNewTokenInfo 添加一个新 poolToken 并写入缓存和文件中
func SetNewTokenInfo(token string, req *tokenModel.TokensRequest) (err error) {
	models := make(map[string]*tokenModel.PoolModelInfo)
	models[common.ChatGPT] = tokenModel.NewPoolModelInfo(req.ChatGPT)
	models[common.Copilot] = tokenModel.NewPoolModelInfo(req.Copilot)
	models[common.CoCopilot] = tokenModel.NewPoolModelInfo(req.CoCopilot)
	models[common.Gemini] = tokenModel.NewPoolModelInfo(req.GeMini)
	tokenInfo := tokenModel.PoolTokenInfo{
		Expire:     0,
		CreateTime: time.Now(),
	}
	tokenJson, err := global.Json.Marshal(tokenInfo)
	if err != nil {
		global.SugarLog.Errorw("SetNewTokenInfo json Marshal error", "err", err, "token", token)
		return
	}
	err = TokenCache.Set(token, tokenJson)
	if err != nil {
		global.SugarLog.Errorw("SetNewTokenInfo set ceche error", "err", err, "token", token)
		return
	}
	//go SetModelTokenCache(token, models)
	// TODO set file
	return
}

// SetModelTokenCache 设置模型Key
//func SetModelTokenCache(token string, models map[string]*tokenModel.PoolModelInfo) {
//	for model, info := range models {
//		// 模型|poolToken
//		indexKey := fmt.Sprintf("%s|%s", model.Name, token)
//		err := TokenCache.Set(indexKey, utils.Int2Byte(0))
//		if err != nil {
//			global.SugarLog.Errorw("SetModelTokenCache set token index error", "indexKey", indexKey, "setVal", utils.Int2Byte(0))
//		}
//		for index, key := range model.Tokens {
//			// 模型|poolToken|index
//			cacheKey := fmt.Sprintf("%s|%s|%d", model.Name, token, index)
//			err := TokenCache.Set(cacheKey, []byte(key))
//			if err != nil {
//				global.SugarLog.Errorw("SetModelTokenCache set token error", "cacheKey", cacheKey, "setVal", key)
//			}
//		}
//	}
//}

// GenerateToken 生成 Token
func GenerateToken(prefix string) string {
	b := make([]byte, 32)
	rand.Read(b)
	// 使用base64编码token
	token := base64.StdEncoding.EncodeToString(b)
	token = fmt.Sprintf("%s-%s", prefix, token)
	return token
}

func CheckTokenParam(req *tokenModel.TokensRequest) string {
	chatGPTLen := len(req.ChatGPT)
	coCopilotLen := len(req.CoCopilot)
	copilotLen := len(req.Copilot)
	geMiniLen := len(req.GeMini)

	if chatGPTLen == 0 && coCopilotLen == 0 && copilotLen == 0 && geMiniLen == 0 {
		return "all tokens are empty"
	}
	if chatGPTLen > 1000 {
		return "ChatGPT key exceeds maximum length"
	}
	if coCopilotLen > 1000 {
		return "CoCopilot key exceeds maximum length"
	}
	if copilotLen > 1000 {
		return "Copilot key exceeds maximum length"
	}
	if geMiniLen > 1000 {
		return "GeMini key exceeds maximum length"
	}
	return ""
}
