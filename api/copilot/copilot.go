package copilot

import (
	"bufio"
	"context"
	"errors"
	"github.com/allegro/bigcache/v3"
	"github.com/dalefeng/chat-api-reverse/global"
	"github.com/dalefeng/chat-api-reverse/model/common/response"
	copilotModel "github.com/dalefeng/chat-api-reverse/model/copilot"
	"github.com/dalefeng/chat-api-reverse/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var CopilotTokenCache *bigcache.BigCache
var TokenExpiredError = errors.New("token expired")

func init() {
	var err error
	CopilotTokenCache, err = bigcache.New(context.Background(), bigcache.DefaultConfig(23*time.Minute))
	if err != nil {
		log.Println("init CopilotTokenCache error ", err.Error())
		panic(err)
	}
	log.Println("init CopilotTokenCache success")
}

type CopilotApi struct {
}

// TokenHandler 从 官方 Copilot 获取到 Github Copilot
func (co *CopilotApi) TokenHandler(c *gin.Context) {
	token, err := utils.GetAuthToken(c, "token")
	if err != nil {
		global.SugarLog.Errorw("TokenHandler get auth token err", "err", err)
		response.FailWithChat(http.StatusUnauthorized, err.Error(), c)
		return
	}
	_, respMap, httpStatus, err := GetCopilotToken(token, false)

	c.Status(httpStatus)
	if err != nil {
		global.SugarLog.Errorw("TokenHandler GetCopilotToken error", "err", err, "token", token)
		c.JSON(httpStatus, respMap)
		return
	}
	c.JSON(httpStatus, respMap)
}

// CoTokenHandler 代理从 CoCopilot 获取到 Github Copilot CoTokenHandler
func (co *CopilotApi) CoTokenHandler(c *gin.Context) {
	token, err := utils.GetAuthToken(c, "token")
	if err != nil {
		global.SugarLog.Errorw("CoTokenHandler get auth token err", "err", err)
		response.FailWithChat(http.StatusUnauthorized, err.Error(), c)
		return
	}
	_, respMap, httpStatus, err := GetCopilotToken(token, true)

	c.Status(httpStatus)
	if err != nil {
		global.SugarLog.Errorw("CoTokenHandler GetCopilotToken error", "err", err, "token", token)
		c.JSON(httpStatus, respMap)
		return
	}
	c.JSON(httpStatus, respMap)
}

// CompletionsHandler 兼容 CoCopilot 的官方 Completions 接口
func (co *CopilotApi) CompletionsHandler(c *gin.Context) {
	var req map[string]interface{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.SugarLog.Errorw("CompletionsHandler bind json err", "err", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	token, err := utils.GetAuthToken(c, "Bearer")
	if err != nil {
		global.SugarLog.Errorw("CompletionsHandler get auth token err", "err", err)
		response.FailWithChat(http.StatusUnauthorized, err.Error(), c)
		return
	}

	// 获取 CopilotToken
	copilotToken, err := GetCopilotTokenWithCache(token)
	if err != nil {
		response.FailWithChat(http.StatusUnauthorized, err.Error(), c)
		return
	}
	err = CompletionsRequest(c, req, copilotToken)
	// 如果 token 过期，重新获取一次 token
	if errors.Is(err, TokenExpiredError) {
		CopilotTokenCache.Delete(token) // 删除缓存
		global.SugarLog.Infow("CompletionsHandler token expired, try get new token", "token", token)
		coCopilotToken, _, _, err := GetCopilotToken(token, true)
		if err != nil {
			global.SugarLog.Errorw("CompletionsHandler http fetch token, Try twice error", "err", err, "token", token)
			response.FailWithChat(http.StatusUnauthorized, err.Error(), c)
			return
		}
		global.SugarLog.Infow("CompletionsHandler http get token is success")
		err = CompletionsRequest(c, req, coCopilotToken)
		if err != nil {
			global.SugarLog.Errorw("CompletionsHandler http fetch token, Try twice error", "err", err, "token", token)
			return
		}
	}
}

// GetCopilotTokenWithCache 先从缓存中获取 CopilotToken，如果缓存中没有，再从 CoCopilot 获取
func GetCopilotTokenWithCache(token string) (copilotToken string, err error) {
	cacheToken, cacheErr := CopilotTokenCache.Get(token)
	if cacheErr != nil {
		global.SugarLog.Infow("CompletionsHandler get cache err, Try http fetch token", "err", cacheErr)
		var tokenErr error
		copilotToken, _, _, tokenErr = GetCopilotToken(token, true)
		if tokenErr != nil {
			global.SugarLog.Errorw("CompletionsHandler http fetch token error", "cacheErr", cacheErr, "tokenErr", tokenErr, "token", token)
			err = tokenErr
			return
		}
	} else {
		copilotToken = string(cacheToken)
		global.SugarLog.Infow("CompletionsHandler get cache success", "copilotToken", copilotToken)
	}
	return
}

// CompletionsRequest 请求 Copilot Completions 接口
func CompletionsRequest(c *gin.Context, req map[string]interface{}, copilotToken string) (err error) {
	url := global.Config.Copilot.CompletionsURL
	client := resty.New()
	resp, err := client.R().
		SetDoNotParseResponse(true).
		SetHeaders(copilotModel.GetCompletionsHeader(copilotToken)).
		SetBody(req).
		Post(url)

	if err != nil {
		global.SugarLog.Errorw("CompletionsRequest http error", "err", err, "url", url, "req", req, "copilotToken", copilotToken)
		response.FailWithChat(http.StatusInternalServerError, err.Error(), c)
		return
	}
	defer resp.RawBody().Close()
	reader := bufio.NewReader(resp.RawBody())

	respContentType := resp.Header().Get("Content-Type")
	if resp.StatusCode() != http.StatusOK {
		global.SugarLog.Warnw("CompletionsRequest respContentType", "respContentType", respContentType)
	}

	if strings.Contains(respContentType, "text/plain") {
		body, err := io.ReadAll(reader)
		if err != nil {
			global.SugarLog.Errorw("CompletionsHandler reader body err", "err", err)
			return err
		}
		bodyStr := strings.TrimRight(string(body), "\n")
		if bodyStr == "unauthorized: token expired" {
			global.SugarLog.Errorw("CompletionsHandler token expired", "body", bodyStr)
			return TokenExpiredError
		}
		global.SugarLog.Infow("CompletionsHandler response error", "body", bodyStr)
		response.FailWithChat(resp.StatusCode(), bodyStr, c)
		return nil
	}

	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream") // 声明数据格式为event stream
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no") // // 禁用nginx缓存,防止nginx会缓存数据导致数据流是一段一段的
	flusher, _ := w.(http.Flusher)

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			global.SugarLog.Errorw("CompletionsHandler reader err", "err", err)
			break
		}
		w.Write(line)
		flusher.Flush()
	}
	return
}

// GetCopilotToken get token from coCopilot
// isCo: true: coCopilot, false: Copilot
func GetCopilotToken(key string, isCo bool) (token string, data map[string]interface{}, httpCode int, err error) {
	data = make(map[string]interface{})
	errData := make(map[string]interface{})
	httpCode = http.StatusInternalServerError
	client := resty.New()
	tokenUrl := global.Config.Copilot.TokenURL
	if isCo {
		tokenUrl = global.Config.Copilot.CoTokenURL
	}
	u, err := url.Parse(tokenUrl)
	if err != nil {
		global.SugarLog.Errorw("GetCopilotToken parse url error", "err", err, "tokenUrl", tokenUrl)
		return
	}
	resp, err := client.R().
		SetHeader("Host", u.Host).
		SetHeader("Authorization", "token "+key).
		SetHeader("Editor-Version", "vscode/1.85.0").
		SetHeader("Editor-Plugin-Version", "copilot-chat/0.11.1").
		SetHeader("User-Agent", "GitHubCopilotChat/0.11.1").
		SetHeader("Accept", "*/*").
		SetResult(&data).
		SetError(&errData).
		Get(tokenUrl)
	if err != nil {
		global.SugarLog.Errorw("GetCopilotToken request http error", "err", err, "url", global.Config.Copilot.CoTokenURL, "key", key)
		return
	}
	httpCode = resp.StatusCode()
	if httpCode != http.StatusOK {
		global.SugarLog.Errorw("GetCopilotToken httpCode!== 200", "statusCode", httpCode, "token", token, "errData", errData)
		data = errData
		return
	}
	t, ok := data["token"]
	if !ok {
		global.SugarLog.Errorw("GetCopilotToken response token is nil", "token", token, "data", data, "errData", errData)
		err = errors.New("response token is nil")
		return
	}
	token = t.(string)
	if token == "" {
		global.SugarLog.Errorw("GetCopilotToken response token is empty", "token", token, "data", data, "errData", errData)
		err = errors.New("response token is empty")
		return
	}
	global.SugarLog.Infow("GetCopilotToken GetCopilotToken Success", "key", key)
	cacheErr := CopilotTokenCache.Set(key, []byte(token))
	if cacheErr != nil {
		global.SugarLog.Errorw("GetCopilotToken set cache err", "err", cacheErr)
	}
	return
}
