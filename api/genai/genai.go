package genai

import (
	"github.com/dalefeng/chat-api-reverse/global"
	"github.com/dalefeng/chat-api-reverse/model/common/response"
	genModel "github.com/dalefeng/chat-api-reverse/model/genai"
	openModel "github.com/dalefeng/chat-api-reverse/model/openai"
	"github.com/dalefeng/chat-api-reverse/pkg/adapter/event"
	genaiAdapter "github.com/dalefeng/chat-api-reverse/pkg/adapter/genai"
	"github.com/dalefeng/chat-api-reverse/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	"io"
	"net/http"
)

type GenApi struct {
}

func (g *GenApi) CompletionsHandler(c *gin.Context) {
	token, err := utils.GetAuthToken(c, "Bearer")
	if err != nil {
		global.SugarLog.Errorw("CompletionsHandler Get hearder token error", "error", err)
		response.FailWithOpenAIError(http.StatusUnauthorized, err.Error(), c)
		return
	}

	var req = &openModel.ChatCompletionRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.FailWithOpenAIError(http.StatusBadRequest, err.Error(), c)
		return
	}

	endpoint := internaloption.WithDefaultEndpoint(global.Config.Gemini.BaseUrl)
	client, err := genai.NewClient(c, option.WithAPIKey(token), endpoint)
	if err != nil {
		global.SugarLog.Errorw("genai.NewClient error", "error", err)
		response.FailWithOpenAIError(http.StatusInternalServerError, err.Error(), c)
	}
	defer client.Close()

	// 适配器
	var gemini genaiAdapter.GenAiModelAdapter
	switch req.Model {
	case genModel.GeminiProVision:
		//gemini = genaiAdapter.NewGeminiProModelAdapter(client)
		panic("not support")
	case openai.GPT432K: // 自定义模型
		gemini = genaiAdapter.NewGeminiProModelAdapter(client)
	default:
		gemini = genaiAdapter.NewGeminiProModelAdapter(client)
	}
	// 不是流式输出
	if !req.Stream {
		panic("not support")
	}

	dataChan, _ := gemini.GenerateStreamContent(c, req)

	c.Stream(func(w io.Writer) bool {
		if data, ok := <-dataChan; ok {
			c.Render(-1, event.Event{Data: "data: " + data})
			return true
		}
		c.Render(-1, event.Event{Data: "data: [DONE]"})
		return false
	})
}
