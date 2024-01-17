package token

import (
	"errors"
	"github.com/dalefengs/chat-api-proxy/global"
	"github.com/dalefengs/chat-api-proxy/model/common"
	"github.com/dalefengs/chat-api-proxy/model/common/response"
	tokenModel "github.com/dalefengs/chat-api-proxy/model/token"
	tokenSvc "github.com/dalefengs/chat-api-proxy/pkg/service/token"
	"github.com/gin-gonic/gin"
)

type TokenApi struct{}

func (a *TokenApi) TokenPoolHandler(c *gin.Context) {
	var req tokenModel.TokensRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.SugarLog.Errorw("TokenPoolHandler param format error", "err", err)
		response.FailWithMessage("param format error", c)
		return
	}
	if msg := tokenSvc.CheckTokenParam(&req); msg != "" {
		response.FailWithMessage(msg, c)
		return
	}
	var poolToken string
	var isChange bool
	var poolTokenInfo *tokenModel.PoolTokenInfo
	if req.Token != "" {
		poolToken = req.Token
		isChange = true
		fileTokenInfo, err := tokenSvc.GetPoolTokenInfoByFile(poolToken)
		// 发生异常且不是 poolToken 不存在异常
		if err != nil && !errors.Is(err, tokenSvc.TokenNotExistsError) {
			global.SugarLog.Errorw("TokenPoolHandler get poolToken error", "err", err)
			response.FailWithMessage("get poolToken error", c)
			return
		}
		if errors.Is(err, tokenSvc.TokenNotExistsError) {
			global.SugarLog.Warnw("TokenPoolHandler Token Not Exists Error", "err", err)
			response.FailWithMessage("token not exists", c)
			return
		}
		isChange = true
		poolTokenInfo, err = tokenSvc.UpdatePoolTokenInfo(poolToken, &req, fileTokenInfo)
	} else {
		poolToken = tokenSvc.GenerateToken("fpk")
		poolTokenInfo, err = tokenSvc.CreatePoolTokenInfo(poolToken, &req)
	}
	if err != nil {
		response.FailWithMessage("save poolToken error", c)
		return
	}
	resp := gin.H{
		"poolToken":   poolToken,
		"copilot":     len(poolTokenInfo.Models[common.Copilot].Tokens),
		"cocopilot":   len(poolTokenInfo.Models[common.CoCopilot].Tokens),
		"chatgpt":     len(poolTokenInfo.Models[common.ChatGPT].Tokens),
		"gemini":      len(poolTokenInfo.Models[common.Gemini].Tokens),
		"create_time": poolTokenInfo.CreateTime.Format("2006-01-02 15:04:05"),
	}
	if isChange {
		resp["update_time"] = poolTokenInfo.UpdateTime.Format("2006-01-02 15:04:05")
	}
	response.OkWithData(resp, c)
}

func (a *TokenApi) TokenInfoHandler(c *gin.Context) {
	poolToken := c.Query("token")
	if poolToken == "" {
		response.FailWithMessage("token is empty", c)
		return
	}
	info, err := tokenSvc.GetPoolTokenInfo(poolToken)
	if err != nil {
		response.FailWithMessage("get token info error", c)
		return
	}
	response.OkWithData(info, c)
}
