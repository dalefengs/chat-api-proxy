package copilot

import (
	"errors"
	"github.com/dalefeng/chat-api-reverse/global"
	"github.com/dalefeng/chat-api-reverse/model/common/response"
	"github.com/dalefeng/chat-api-reverse/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type CopilotApi struct {
}

func (co *CopilotApi) Token(c *gin.Context) {
	token, err := utils.GetAuthToken(c)
	if err != nil {
		global.SugarLog.Errorw("get token err", "err", err)
		response.FailWithChat(http.StatusUnauthorized, err.Error(), c)
		return
	}
	_, respMap, httpStatus, err := GetCoCopilotToken(token)
	c.Status(httpStatus)
	if err != nil {
		global.SugarLog.Errorw("GetCoCopilotToken", "err", err, "token", token)
		c.JSON(httpStatus, respMap)
		return
	}
	c.JSON(httpStatus, respMap)
}

func (co *CopilotApi) Completions(c *gin.Context) {
	response.OkWithMessage("success", c)
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
	return
}
