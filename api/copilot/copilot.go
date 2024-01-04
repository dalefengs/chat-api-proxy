package copilot

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/dalefeng/chat-api-reverse/global"
	"github.com/dalefeng/chat-api-reverse/model/common/response"
	copilotModel "github.com/dalefeng/chat-api-reverse/model/copilot"
	"github.com/dalefeng/chat-api-reverse/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var CopilotTokenCache *bigcache.BigCache

func init() {
	var err error
	CopilotTokenCache, err = bigcache.New(context.Background(), bigcache.DefaultConfig(22*time.Minute))
	if err != nil {
		log.Println("init CopilotTokenCache error ", err.Error())
		panic(err)
	}
	log.Println("init CopilotTokenCache success")
}

type CopilotApi struct {
}

// CoToken 代理从 CoCopilot 获取到 Github Copilot CoToken
func (co *CopilotApi) CoToken(c *gin.Context) {
	token, err := utils.GetAuthToken(c, "token")
	if err != nil {
		global.SugarLog.Errorw("get auth token err", "err", err)
		response.FailWithChat(http.StatusUnauthorized, err.Error(), c)
		return
	}
	copilotToken, respMap, httpStatus, err := GetCoCopilotToken(token)
	cacheErr := CopilotTokenCache.Set(token, []byte(copilotToken))
	if cacheErr != nil {
		global.SugarLog.Errorw("CoToken set cache err", "err", cacheErr)
	}
	c.Status(httpStatus)
	if err != nil {
		global.SugarLog.Errorw("GetCoCopilotToken", "err", err, "token", token)
		c.JSON(httpStatus, respMap)
		return
	}
	c.JSON(httpStatus, respMap)
}

// CoCompletions 兼容 CoCopilot 的官方 Completions 接口
func (co *CopilotApi) CoCompletions(c *gin.Context) {
	var req map[string]interface{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.SugarLog.Errorw("CoCompletions bind json err", "err", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	token, err := utils.GetAuthToken(c, "Bearer")
	if err != nil {
		global.SugarLog.Errorw("get auth token err", "err", err)
		response.FailWithChat(http.StatusUnauthorized, err.Error(), c)
		return
	}

	url := global.Config.Copilot.CompletionsURL
	client := resty.New()
	resp, err := client.R().
		SetDoNotParseResponse(true).
		SetHeaders(copilotModel.GetCompletionsHeader(token)).
		SetAuthToken(token).
		SetBody(req).
		Post(url)

	if err != nil {
		global.SugarLog.Errorw("request http error", "err", err, "url", url, "req", req, "token", token)
		response.FailWithChat(http.StatusInternalServerError, err.Error(), c)
		return
	}
	defer resp.RawBody().Close()
	reader := bufio.NewReader(resp.RawBody())

	respContentType := resp.Header().Get("Content-Type")
	if !strings.Contains(respContentType, "text/event-stream") {
		body, err := io.ReadAll(reader)
		if err != nil {
			global.SugarLog.Errorw("reader body err", "err", err)
			return
		}
		var data map[string]interface{}
		jsonErr := jsoniter.Unmarshal(body, &data)
		if jsonErr != nil {
			global.SugarLog.Errorw("CoCompletions response json unmarshal error", "err", jsonErr, "body", string(body))
			response.FailWithChat(http.StatusInternalServerError, jsonErr.Error(), c)
			return
		}
		c.JSON(resp.StatusCode(), data)
		return
	}

	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream") // 声明数据格式为event stream
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no") // // 禁用nginx缓存,防止nginx会缓存数据导致数据流是一段一段的
	flusher, _ := w.(http.Flusher)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			global.SugarLog.Errorw("reader err", "err", err)
			break
		}
		if line == "\n" {
			continue
		}
		fmt.Fprintf(w, line+"\n")
		flusher.Flush()
	}
	flusher.Flush()
}

// GetCoCopilotToken get token from coCopilot
func GetCoCopilotToken(key string) (token string, data map[string]interface{}, httpCode int, err error) {
	data = make(map[string]interface{})
	errData := make(map[string]interface{})
	httpCode = http.StatusInternalServerError
	client := resty.New()
	resp, err := client.R().
		SetHeader("Host", "api.cocopilot.org").
		SetHeader("Authorization", "token "+key).
		SetHeader("Editor-Version", "vscode/1.85.0").
		SetHeader("Editor-Plugin-Version", "copilot-chat/0.11.1").
		SetHeader("User-Agent", "GitHubCopilotChat/0.11.1").
		SetHeader("Accept", "*/*").
		SetResult(&data).
		SetError(&errData).
		Get(global.Config.Copilot.CoTokenURL)
	if err != nil {
		global.SugarLog.Errorw("request http error", "err", err, "url", global.Config.Copilot.CoTokenURL, "key", key)
		return
	}
	httpCode = resp.StatusCode()
	if httpCode != http.StatusOK {
		global.SugarLog.Errorw("httpCode!== 200", "statusCode", httpCode, "token", token, "errData", errData)
		data = errData
		return
	}
	t, ok := data["token"]
	if !ok {
		global.SugarLog.Errorw("response token is nil", "token", token, "data", data, "errData", errData)
		err = errors.New("response token is nil")
		return
	}
	token = t.(string)
	if token == "" {
		global.SugarLog.Errorw("response token is empty", "token", token, "data", data, "errData", errData)
		err = errors.New("response token is empty")
		return
	}
	global.SugarLog.Infow("GetCoCopilotToken", "key", "token", token)
	return
}
