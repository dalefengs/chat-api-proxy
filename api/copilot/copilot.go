package copilot

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/dalefengs/chat-api-proxy/api/genai"
	"github.com/dalefengs/chat-api-proxy/global"
	"github.com/dalefengs/chat-api-proxy/model/common/response"
	"github.com/dalefengs/chat-api-proxy/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/sashabaranov/go-openai"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var TokenCache *freecache.Cache
var TokenExpiredError = errors.New("token expired")
var Client = resty.New()

type StatisticsType string

const (
	TypeAuthToken   StatisticsType = "authToken"
	TypeCacheToken  StatisticsType = "cacheToken"
	TypeCompletions StatisticsType = "completions"
)

var AuthCount = make(map[string]map[StatisticsType]int)

func init() {
	cacheSize := 10 * 1024 * 1024 // 10M
	TokenCache = freecache.NewCache(cacheSize)
	log.Println("init CopilotTokenCache success")
	go ClearAuthCount()
}

type CopilotApi struct {
}

// TokenHandler 从 官方 Copilot 获取到 Github Copilot
func (co *CopilotApi) TokenHandler(c *gin.Context) {
	token, err := utils.GetAuthToken(c, "token")
	if err != nil {
		global.SugarLog.Errorw("TokenHandler get auth token err", "err", err)
		response.FailWithOpenAIError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	day := time.Now().Format("2006-01-02")
	if dayCount, ok := AuthCount[day]; ok {
		if authCount, ok := dayCount[TypeAuthToken]; ok && authCount > 200 {
			global.SugarLog.Errorw("TokenHandler get auth token err,  authorization limit > 200")
			response.FailWithOpenAIError(http.StatusUnauthorized, "exceeded authorization limit", c)
			return
		}
	}

	tokenRawCache, err := GetTokenRawInfoCache(token)
	if err == nil && len(tokenRawCache) > 0 {
		global.SugarLog.Infow("TokenHandler get token raw cache", "token", token)
		c.JSON(http.StatusOK, tokenRawCache)
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
		response.FailWithOpenAIError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	tokenRawCache, err := GetTokenRawInfoCache(token)
	if err == nil && len(tokenRawCache) > 0 {
		global.SugarLog.Infow("CoTokenHandler get token raw cache", "token", token)
		c.JSON(http.StatusOK, tokenRawCache)
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

// CompletionsHandler 兼容 CoCopilot 的官方 CompletionsHandler 接口
func (co *CopilotApi) CompletionsHandler(c *gin.Context) {
	var req map[string]interface{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.SugarLog.Errorw("CompletionsHandler bind json err", "err", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	if model, ok := req["model"]; ok && model == "gemini-pro" || model == openai.GPT432K {
		global.SugarLog.Debugw("CompletionsHandler gemini-pro model", "model", model)
		genApi := &genai.GenApi{}
		genApi.CompletionsHandler(c)
		return
	}

	token, err := utils.GetAuthToken(c, "Bearer")
	if err != nil {
		global.SugarLog.Errorw("CompletionsHandler get auth token err", "err", err.Error())
		response.FailWithOpenAIError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	// 获取 CopilotToken
	copilotToken, err := GetCopilotTokenWithCache(token)
	if err != nil {
		response.FailWithOpenAIError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	IncAuthCount(TypeCompletions)
	err = CompletionsRequest(c, req, copilotToken)
	// 如果 token 过期，重新获取一次 token
	if errors.Is(err, TokenExpiredError) {
		TokenCache.Del([]byte(token)) // 删除缓存
		global.SugarLog.Infow("CompletionsHandler token expired, try get new token", "token", token)
		coCopilotToken, _, _, coErr := GetCopilotToken(token, true)
		if coErr != nil {
			global.SugarLog.Errorw("CompletionsHandler http fetch token, Try twice error", "err", coErr, "token", token)
			response.FailWithOpenAIError(http.StatusUnauthorized, coErr.Error(), c)
			return
		}
		global.SugarLog.Infow("CompletionsHandler http get token is success")
		err = CompletionsRequest(c, req, coCopilotToken)
		if err != nil {
			global.SugarLog.Errorw("CompletionsHandler CompletionsRequest retry request error", "err", err)
			response.FailWithOpenAIError(http.StatusBadGateway, err.Error(), c)
			return
		}
	} else if err != nil {
		global.SugarLog.Warnw("CompletionsHandler CompletionsRequest request error", "err", err)
		response.FailWithOpenAIError(http.StatusInternalServerError, err.Error(), c)
		return
	}
}

// CompletionsOfficialHandler Copilot 的官方 CompletionsHandler 接口，使用官方token
func (co *CopilotApi) CompletionsOfficialHandler(c *gin.Context) {
	var req map[string]interface{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.SugarLog.Errorw("CompletionsOfficialHandler bind json err", "err", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	token, err := utils.GetAuthToken(c, "Bearer")
	if err != nil {
		global.SugarLog.Errorw("CompletionsOfficialHandler get auth token err", "err", err.Error())
		response.FailWithOpenAIError(http.StatusUnauthorized, err.Error(), c)
		return
	}
	err = CompletionsRequest(c, req, token)
	if err != nil {
		global.SugarLog.Warnw("CompletionsOfficialHandler CompletionsRequest request error", "err", err, "token", token)
		response.FailWithOpenAIError(http.StatusInternalServerError, err.Error(), c)
		return
	}
}

func (co *CopilotApi) CountHandler(c *gin.Context) {
	response.OkWithData(AuthCount, c)
}

// GetCopilotTokenWithCache 先从缓存中获取 CopilotToken，如果缓存中没有，再从 CoCopilot 获取
func GetCopilotTokenWithCache(token string) (copilotToken string, err error) {
	cacheToken, cacheErr := TokenCache.Get([]byte(token))
	if cacheErr != nil {
		global.SugarLog.Infow("CompletionsHandler get cache err, Try http fetch token", "err", cacheErr.Error(), "token", token)
		var tokenErr error
		copilotToken, _, _, tokenErr = GetCopilotToken(token, true)
		if tokenErr != nil {
			global.SugarLog.Errorw("CompletionsHandler http fetch token error", "cacheErr", cacheErr, "tokenErr", tokenErr, "token", token)
			err = tokenErr
			return
		}
	} else {
		IncAuthCount(TypeCacheToken)
		copilotToken = string(cacheToken)
		global.SugarLog.Infow("CompletionsHandler get cache success")
	}
	return
}

// GetTokenRawInfoCache 获取 token 的原始信息缓存
func GetTokenRawInfoCache(token string) (respMap map[string]any, err error) {
	tokenRawCache, err := TokenCache.Get([]byte(token + "_raw"))
	if err != nil {
		return
	}
	if tokenRawCache == nil {
		global.SugarLog.Infow("GetTokenRawInfoCache token raw cache is nil", "token", token)
		return
	}
	err = jsoniter.Unmarshal(tokenRawCache, &respMap)
	if err != nil {
		global.SugarLog.Errorw("GetTokenRawInfoCache get token raw cache json unmarshal error", "err", err)
		return
	}
	return
}

// CompletionsRequest 请求 Copilot CompletionsHandler 接口
func CompletionsRequest(c *gin.Context, req map[string]interface{}, copilotToken string) (err error) {
	completionsURL := global.Config.Copilot.CompletionsURL
	resp, err := Client.SetRetryCount(1).R().
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if err != nil && strings.Contains(err.Error(), "connection reset by peer") {
				IncAuthCount(TypeCompletions)
				global.SugarLog.Warnw("CompletionsRequest Client connection reset by peer", "err", err.Error())
			}
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				global.SugarLog.Warnw("CompletionsRequest Client timeout err, retry", "err", netErr.Error())
				return true
			}
			return false
		}).
		SetDoNotParseResponse(true).
		SetHeaders(GetCompletionsHeader(copilotToken)).
		SetBody(req).
		Post(completionsURL)

	if err != nil {
		global.SugarLog.Errorw("CompletionsRequest http error", "err", err, "completionsURL", completionsURL, "req", req, "copilotToken", copilotToken)
		response.FailWithOpenAIError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	defer resp.RawBody().Close()
	reader := bufio.NewReader(resp.RawBody())

	respContentType := resp.Header().Get("Content-Type")
	if resp.StatusCode() != http.StatusOK {
		global.SugarLog.Warnw("CompletionsRequest respContentType", "respContentType", respContentType, "statusCode", resp.StatusCode())
	}

	w := c.Writer
	global.SugarLog.Infow("CompletionsRequest start of processor", "respContentType", respContentType, "statusCode", resp.StatusCode())
	if strings.Contains(respContentType, "text/plain") { // 有错误信息
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
		global.SugarLog.Infow("CompletionsHandler response error", "body", bodyStr, "copilotToken", copilotToken)
		response.FailWithOpenAIError(resp.StatusCode(), bodyStr, c)
		return nil
	}
	if b, ok := req["stream"]; ok && b.(bool) == false { // json 格式 非流式
		utils.SetJsonHeaders(c)
		flusher, _ := w.(http.Flusher)
		body, readErr := io.ReadAll(reader)
		if readErr != nil {
			global.SugarLog.Errorw("CompletionsHandler reader body err", "respContentType", respContentType, "req", req, "err", readErr)
			return readErr
		}
		w.Write(body)
		flusher.Flush()
		global.SugarLog.Infow("CompletionsRequest end of processor, stream is true")
		return
	}

	utils.SetEventStreamHeaders(c)
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
	global.SugarLog.Infow("CompletionsRequest end of processor")
	return
}

// GetCopilotToken get token from coCopilot
// isCo: true: coCopilot, false: Copilot
func GetCopilotToken(key string, isCo bool) (token string, data map[string]interface{}, httpCode int, err error) {
	data = make(map[string]interface{})
	errDataMap := make(map[string]interface{})
	httpCode = http.StatusInternalServerError
	tokenUrl := global.Config.Copilot.TokenURL
	if isCo {
		tokenUrl = global.Config.Copilot.CoTokenURL
	}
	u, err := url.Parse(tokenUrl)
	if err != nil {
		global.SugarLog.Errorw("GetCopilotToken parse url error", "err", err, "tokenUrl", tokenUrl)
		return
	}
	IncAuthCount(TypeAuthToken)
	resp, err := Client.SetRetryCount(1).R().
		AddRetryCondition(func(r *resty.Response, err error) bool {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				IncAuthCount(TypeAuthToken)
				global.SugarLog.Warnw("GetCopilotToken Client timeout err, retry", "err", netErr.Error())
				return true
			}
			return false
		}).
		SetHeader("Host", u.Host).
		SetHeader("Authorization", "token "+key).
		SetHeader("Editor-Version", "vscode/1.85.0").
		SetHeader("Editor-Plugin-Version", "copilot-chat/0.11.1").
		SetHeader("User-Agent", "GitHubCopilotChat/0.11.1").
		SetHeader("Accept", "*/*").
		SetResult(&data).
		SetError(&errDataMap).
		Get(tokenUrl)

	go func() {
		jsonData, _ := jsoniter.MarshalIndent(AuthCount, "", "  ")
		global.SugarLog.Infof("request statistics \n %s", string(jsonData))
	}()
	if err != nil {
		global.SugarLog.Errorw("GetCopilotToken request http error", "err", err, "url", global.Config.Copilot.CoTokenURL, "key", key)
		return
	}
	httpCode = resp.StatusCode()
	if httpCode != http.StatusOK {
		global.SugarLog.Errorw("GetCopilotToken httpCode!== 200", "statusCode", httpCode, "key", key, "errDataMap", errDataMap)
		data = errDataMap
		msg := strings.Builder{}
		if message, ok := errDataMap["message"]; ok && message != nil {
			msg.WriteString("message: ")
			msg.WriteString(message.(string))
		}
		if detail, ok := errDataMap["error_details"]; ok && detail != nil {
			detailMap := detail.(map[string]any)
			message := detailMap["message"]
			if message != nil {
				msg.WriteString(" error_message: ")
				msg.WriteString(message.(string))
			}
		}
		err = fmt.Errorf("get copilot token error, %s", msg.String())
		return
	}
	t, ok := data["token"]
	if !ok {
		global.SugarLog.Errorw("GetCopilotToken response token is nil", "token", token, "data", data, "errDataMap", errDataMap)
		err = errors.New("response token is nil")
		return
	}
	token = t.(string)
	if token == "" {
		global.SugarLog.Errorw("GetCopilotToken response token is empty", "token", token, "data", data, "errDataMap", errDataMap)
		err = errors.New("response token is empty")
		return
	}

	expires := 700
	expiresAt, ok := data["expires_at"]
	var expiresTime string
	if ok {
		switch expiresAt.(type) {
		case int, int64:
			expiresAtInt := expiresAt.(int64)
			if expiresAtInt > 0 {
				expires = int(expiresAtInt - time.Now().Unix())
				expiresTime = time.Unix(expiresAtInt, 0).Format("2006-01-02 15:04:05")
			}
		case float64:
			expiresAtFloat := expiresAt.(float64)
			if expiresAtFloat > 0 {
				expires = int(int64(expiresAtFloat) - time.Now().Unix())
				expiresTime = time.Unix(int64(expiresAtFloat), 0).Format("2006-01-02 15:04:05")
			}
		default:
			expiresTime = "expiresAt type error"
			global.SugarLog.Errorw("GetCopilotToken expires_at type error", "expiresAt", expiresAt)
		}
	}

	global.SugarLog.Infow("GetCopilotToken HTTP Authorisation token success", "key", key, "expires", expires, "expiresTime", expiresTime)

	cacheErr := TokenCache.Set([]byte(key), []byte(token), expires)
	if cacheErr != nil {
		global.SugarLog.Errorw("GetCopilotToken set token cache err", "err", cacheErr)
	}
	rawTokenInfo, err := jsoniter.Marshal(data)
	if err != nil {
		global.SugarLog.Errorw("GetCopilotToken token raw info json marshal error", "err", err)
		return
	}
	cacheErr = TokenCache.Set([]byte(key+"_raw"), rawTokenInfo, expires)
	if cacheErr != nil {
		global.SugarLog.Errorw("GetCopilotToken set token raw cache err", "err", cacheErr)
	}
	return
}

// GetCompletionsHeader 获取 Copilot CompletionsHandler 接口的 Header
func GetCompletionsHeader(token string) map[string]string {
	uid := uuid.New().String()
	headersMap := map[string]string{
		"Host":                        "api.githubcopilot.com",
		"Accept-Encoding":             "gzip, deflate, br",
		"Accept":                      "*/*",
		"Authorization":               "Bearer " + token,
		"X-RequestType-Id":            uid,
		"X-Github-CopilotApi-Version": "2023-07-07",
		"Vscode-Sessionid":            uid + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10),
		"vscode-machineid":            utils.GenHexStr(64),
		"Editor-Version":              "vscode/1.85.0",
		"Editor-Plugin-Version":       "copilot-chat/0.11.1",
		"Openai-Organization":         "github-copilot",
		"Copilot-Integration-Id":      "vscode-chat",
		"Openai-Intent":               "conversation-panel",
		"User-Agent":                  "GitHubCopilotChat/0.11.1",
	}
	return headersMap
}

func IncAuthCount(authType StatisticsType) {
	day := time.Now().Format("2006-01-02")
	if _, ok := AuthCount[day]; !ok {
		AuthCount[day] = make(map[StatisticsType]int)
		AuthCount[day][authType] = 0
	}
	AuthCount[day][authType]++
}

// ClearAuthCount 每天 0 点清空7天前统计数据
func ClearAuthCount() {
	t := time.Now()
	next := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
	d := next.Sub(t)

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t := time.Now()
			date := t.AddDate(0, 0, -7).Format("2006-01-02")
			if _, ok := AuthCount[date]; ok {
				delete(AuthCount, date)
			}
			global.SugarLog.Infow("ClearAuthCount success", "clearDay", date)
			ticker.Reset(24 * time.Hour)
		}
	}
}

var pacContent string
var pacTpl = `function FindProxyForURL(url, host) {
	if (host == 'githubusercontent.com' || host == 'githubcopilot.com' || host == 'cocopilot.net' || host == 'google.com') {
		return '%s'
	}
	return 'DIRECT'
}
`
var directPac = `function FindProxyForURL(url, host) {
	return 'DIRECT'
}
`

func (co *CopilotApi) ProxyPacHandler(ctx *gin.Context) {
	proxy := global.Config.Copilot.Proxy.HTTP
	if pacContent == "" && proxy != "" {
		pacContent = fmt.Sprintf(pacTpl, proxy)
	}
	if pacContent == "" {
		pacContent = directPac
	}
	if proxy != "" && pacContent == directPac {
		pacContent = fmt.Sprintf(pacTpl, proxy)
	}
	ctx.Header("Content-Type", "application/x-ns-proxy-autoconfig")
	ctx.String(http.StatusOK, pacContent)
}
