package openai

import (
	"bufio"
	"fmt"
	"github.com/dalefeng/chat-api-reverse/global"
	"github.com/dalefeng/chat-api-reverse/model/common/response"
	"github.com/dalefeng/chat-api-reverse/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"strings"
)

type OpenApi struct {
}

func (co *OpenApi) Completions(c *gin.Context) {
	var req map[string]interface{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.SugarLog.Errorw("Completions bind json err", "err", err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	token, err := utils.GetAuthToken(c, "Bearer")
	if err != nil {
		global.SugarLog.Errorw("get token err", "err", err)
		response.FailWithOpenAIError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	url := global.Config.OpenAi.BaseURL + "/v1/chat/completions"
	client := resty.New()
	resp, err := client.R().
		SetDoNotParseResponse(true).
		SetHeader("Content-Type", "application/json").
		SetAuthToken(token).
		SetBody(req).
		Post(url)

	if err != nil {
		global.SugarLog.Errorw("request http error", "err", err, "url", url, "req", req, "token", token)
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
		jsonErr := global.Json.Unmarshal(body, &data)
		if jsonErr != nil {
			response.FailWithOpenAIError(resp.StatusCode(), jsonErr.Error(), c)
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
}
