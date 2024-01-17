package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/dalefengs/chat-api-proxy/global"
	"github.com/dalefengs/chat-api-proxy/model/common"
	tokenModel "github.com/dalefengs/chat-api-proxy/model/token"
	"github.com/dalefengs/chat-api-proxy/pkg/cache"
	"github.com/dalefengs/chat-api-proxy/utils"
	"time"
)

var TokenNotExistsError = fmt.Errorf("token not exists")

// GetPoolTokenInfo 获取PoolToken详情信息
func GetPoolTokenInfo(poolToken string) (tokenInfo *tokenModel.PoolTokenInfo, err error) {
	tokenCache, err := cache.TokenCache.Get([]byte(poolToken))
	if err == nil {
		err = global.Json.Unmarshal(tokenCache, &tokenInfo)
		if err != nil {
			global.SugarLog.Errorw("GetPoolTokenInfo cache json Unmarshal error", "err", err, "poolToken", poolToken)
			return
		}
		// 获取模型详情
		tokenInfo.Models = make(map[string]*tokenModel.PoolModelInfo)
		models := []string{common.ChatGPT, common.Copilot, common.CoCopilot, common.Gemini}
		for _, modelName := range models {
			poolModelInfoBytes, cacheErr := cache.TokenCache.Get([]byte(fmt.Sprintf("%s|%s", modelName, poolToken)))
			if cacheErr != nil {
				global.SugarLog.Errorw("GetPoolTokenInfo get pool model info cache error", "err", cacheErr, "poolToken", poolToken, "modelName", modelName)
				continue
			}
			var poolModelInfo tokenModel.PoolModelInfo
			jsonErr := global.Json.Unmarshal(poolModelInfoBytes, &poolModelInfo)
			if jsonErr != nil {
				global.SugarLog.Errorw("GetPoolTokenInfo poolModelInfo json Unmarshal error", "err", jsonErr, "poolToken", poolToken, "modelName", modelName)
				continue
			}
			tokenInfo.Models[modelName] = &poolModelInfo
		}
		return
	}
	if errors.Is(err, freecache.ErrNotFound) {
		global.SugarLog.Warnw("GetPoolTokenInfoByCache poolToken not found", "poolToken", poolToken)
	} else if err != nil {
		global.SugarLog.Errorw("GetPoolTokenInfoByCache poolToken cache error", "err", err, "poolToken", poolToken)
	}
	// 从文件中获取
	tokenInfo, err = GetPoolTokenInfoByFile(poolToken)
	if err != nil {
		return
	}
	global.SugarLog.Debugw("从文件中获取到 Token, 写入缓存。", "poolToken", poolToken)
	go SetModelTokenCache(poolToken, tokenInfo)
	return
}

// GetPoolTokenInfoByFile 从文件中获取 poolToken 信息
func GetPoolTokenInfoByFile(token string) (tokenInfo *tokenModel.PoolTokenInfo, err error) {
	path := utils.GetTokenCacheFilePath(token)
	exists, err := utils.PathExists(path)
	if err != nil {
		return
	}
	if !exists {
		err = TokenNotExistsError
		return
	}
	fileData, err := utils.ReadGzipFile(path)
	if err != nil {
		return
	}
	if fileData == nil {
		err = TokenNotExistsError
		return
	}
	err = global.Json.Unmarshal(fileData, &tokenInfo)
	if err != nil {
		global.SugarLog.Errorw("GetPoolTokenInfoByFile json Unmarshal error", "err", err, "token", token)
		return
	}
	return
}

// UpdatePoolTokenInfo 更新 poolToken 信息
func UpdatePoolTokenInfo(poolToken string, req *tokenModel.TokensRequest, fileTokenInfo *tokenModel.PoolTokenInfo) (tokenInfo *tokenModel.PoolTokenInfo, err error) {
	reqModelMap := make(map[string]map[string]struct{})
	for _, token := range req.ChatGPT {
		if reqModelMap[common.ChatGPT] == nil {
			reqModelMap[common.ChatGPT] = make(map[string]struct{})
		}
		reqModelMap[common.ChatGPT][token] = struct{}{}
	}
	for _, token := range req.Copilot {
		if reqModelMap[common.Copilot] == nil {
			reqModelMap[common.Copilot] = make(map[string]struct{})
		}
		reqModelMap[common.Copilot][token] = struct{}{}
	}
	for _, token := range req.CoCopilot {
		if reqModelMap[common.CoCopilot] == nil {
			reqModelMap[common.CoCopilot] = make(map[string]struct{})
		}
		reqModelMap[common.CoCopilot][token] = struct{}{}
	}
	for _, token := range req.Gemini {
		if reqModelMap[common.Gemini] == nil {
			reqModelMap[common.Gemini] = make(map[string]struct{})
		}
		reqModelMap[common.Gemini][token] = struct{}{}
	}
	now := time.Now()

	for modelName, modelInfo := range fileTokenInfo.Models {
		newTokens := make([]tokenModel.TokenInfo, 0, len(modelInfo.Tokens))
		modelMap := make(map[string]struct{}, len(modelInfo.Tokens))
		for _, tokenData := range modelInfo.Tokens {
			modelMap[tokenData.Token] = struct{}{}
		}
		// 缓存中存在，请求中不存在的token需要剔除
		for _, tokenData := range modelInfo.Tokens {
			if _, ok := reqModelMap[modelName][tokenData.Token]; ok {
				newTokens = append(newTokens, tokenData)
				continue
			}
		}
		// 缓存中不存在，请求中存在的token需要添加缓存
		for reqToken, _ := range reqModelMap[modelName] {
			if _, ok := modelMap[reqToken]; !ok {
				newTokens = append(newTokens, tokenModel.TokenInfo{
					Token:  reqToken,
					Expire: 0,
				})
			}
		}
		modelInfo.Tokens = newTokens
		modelInfo.Count = len(newTokens)
	}
	fileTokenInfo.Token = poolToken
	fileTokenInfo.UpdateTime = &now
	err = SetTokenTokenInfoFileCache(poolToken, fileTokenInfo)
	if err != nil {
		return
	}
	tokenInfo = fileTokenInfo
	// 清理空模型缓存
	go ClearEmptyModelTokenCache(poolToken, req)
	go SetModelTokenCache(poolToken, fileTokenInfo)
	return
}

// CreatePoolTokenInfo 添加一个新 poolToken 并写入缓存和文件中
func CreatePoolTokenInfo(token string, req *tokenModel.TokensRequest) (tokenInfo *tokenModel.PoolTokenInfo, err error) {
	models := make(map[string]*tokenModel.PoolModelInfo)
	models[common.ChatGPT] = tokenModel.NewPoolModelInfo(req.ChatGPT)
	models[common.Copilot] = tokenModel.NewPoolModelInfo(req.Copilot)
	models[common.CoCopilot] = tokenModel.NewPoolModelInfo(req.CoCopilot)
	models[common.Gemini] = tokenModel.NewPoolModelInfo(req.Gemini)

	tokenInfo = &tokenModel.PoolTokenInfo{
		Token:      token,
		CreateTime: time.Now(),
	}
	// 写入文件缓存
	tokenInfo.Models = models
	err = SetTokenTokenInfoFileCache(token, tokenInfo)
	if err != nil {
		global.SugarLog.Errorw("CreatePoolTokenInfo set cache file error", "err", err, "token", token)
		return
	}
	// 写入模型 Token 缓存
	go SetModelTokenCache(token, tokenInfo)
	return
}

// SetTokenTokenInfoFileCache 设置 token 文件缓存
func SetTokenTokenInfoFileCache(token string, poolTokenInfo *tokenModel.PoolTokenInfo) error {
	path := fmt.Sprintf("%s/%s", global.UserHomeCacheDir, token)
	return utils.ZipMarshalAndWriteToFile(path, poolTokenInfo)
}

// SetModelTokenCache 设置模型 Token 缓存
func SetModelTokenCache(poolToken string, tokenInfo *tokenModel.PoolTokenInfo) {
	global.SugarLog.Debugw("SetModelTokenCache start set cache", "poolToken", poolToken)

	var tokenNotWithModel tokenModel.PoolTokenInfo
	err := utils.CopyStruct(tokenInfo, &tokenNotWithModel)
	if err != nil {
		global.SugarLog.Errorw("GetPoolTokenInfo CopyStruct error", "error", err)
	}
	tokenNotWithModel.Models = make(map[string]*tokenModel.PoolModelInfo)
	tokenNotWithModelJson, _ := global.Json.Marshal(tokenNotWithModel)
	err = cache.TokenCache.Set([]byte(poolToken), tokenNotWithModelJson, cache.PoolTokenExpired)
	if err != nil {
		global.SugarLog.Errorw("CreatePoolTokenInfo set cache error", "err", err, "token", poolToken)
		return
	}
	for modelName, modelInfo := range tokenInfo.Models {
		// 模型|poolToken
		// value: index
		poolModelInfoKey := fmt.Sprintf("%s|%s", modelName, poolToken)
		modelInfoJson, err := global.Json.Marshal(modelInfo)
		err = cache.TokenCache.Set([]byte(poolModelInfoKey), modelInfoJson, cache.PoolTokenExpired)
		if err != nil {
			global.SugarLog.Errorw("SetModelTokenCache set modelInfoJson error", "poolModelInfoKey", poolModelInfoKey, "setVal")
		}
	}
	global.SugarLog.Debugw("SetModelTokenCache end", "poolToken", poolToken)
}

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
	geMiniLen := len(req.Gemini)

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
		return "Gemini key exceeds maximum length"
	}
	return ""
}

// ClearEmptyModelTokenCache 清理空的模型 Token 缓存
func ClearEmptyModelTokenCache(poolToken string, req *tokenModel.TokensRequest) {
	clearKey := make([]string, 0, 4)
	if len(req.ChatGPT) == 0 {
		clearKey = append(clearKey, fmt.Sprintf("%s|%s", common.ChatGPT, poolToken))
	}
	if len(req.CoCopilot) == 0 {
		clearKey = append(clearKey, fmt.Sprintf("%s|%s", common.CoCopilot, poolToken))
	}
	if len(req.Copilot) == 0 {
		clearKey = append(clearKey, fmt.Sprintf("%s|%s", common.Copilot, poolToken))
	}
	if len(req.Gemini) == 0 {
		clearKey = append(clearKey, fmt.Sprintf("%s|%s", common.Gemini, poolToken))
	}
	for _, key := range clearKey {
		cache.TokenCache.Del([]byte(key))
	}
	return
}
