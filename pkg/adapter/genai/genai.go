package genai

import (
	"errors"
	"fmt"
	"github.com/dalefeng/chat-api-reverse/global"
	genModel "github.com/dalefeng/chat-api-reverse/model/genai"
	openModel "github.com/dalefeng/chat-api-reverse/model/openai"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"time"
)

type GenAiModelAdapter interface {
	GenerateStreamContent(ctx *gin.Context, req *openModel.ChatCompletionRequest) (<-chan string, error) // 流式生成内容
}

type GeminiProModelAdapter struct {
	client *genai.Client
}

// NewGeminiProModelAdapter 创建 GeminiProModelAdapter
func NewGeminiProModelAdapter(client *genai.Client) *GeminiProModelAdapter {
	return &GeminiProModelAdapter{client: client}
}

// GenerateStreamContent 流式生成内容
func (g GeminiProModelAdapter) GenerateStreamContent(c *gin.Context, req *openModel.ChatCompletionRequest) (<-chan string, error) {
	generativeModel := g.client.GenerativeModel(genModel.GeminiPro)
	setGenAiModelByOpenaiRequest(generativeModel, req)
	// Initialize the chat
	cs := generativeModel.StartChat()
	// 设置历史消息
	setGenAiChatHistoryByOpenaiRequest(cs, req)
	// 用户发送的最新一条信息
	prompt := genai.Text(req.Messages[len(req.Messages)-1].StringContent())
	iter := cs.SendMessageStream(c, prompt)

	dataChan := make(chan string)
	// 处理流式输出迭代器
	go HandleStreamIter(iter, dataChan)

	return dataChan, nil
}

func HandleStreamIter(genIter *genai.GenerateContentResponseIterator, dataChan chan string) {
	defer close(dataChan)
	respID := uuid.New().String()
	created := time.Now().Unix()
	for {
		genAiResp, err := genIter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			global.SugarLog.Errorw("SendStream iter Error", "error", err, "models", genModel.GeminiPro)
			apiErr := openai.APIError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
			resp, _ := global.Json.Marshal(apiErr)
			dataChan <- string(resp)
			break
		}
		openaiResp := genAiResponseToStreamCompletionResponse(genAiResp, respID, created)
		resp, _ := global.Json.Marshal(openaiResp)
		dataChan <- string(resp)
	}
}

// 设置 genai 模型参数
func setGenAiModelByOpenaiRequest(model *genai.GenerativeModel, req *openModel.ChatCompletionRequest) {
	if req.MaxTokens != 0 {
		model.MaxOutputTokens = &req.MaxTokens
	}
	if req.Temperature != 0 {
		model.Temperature = &req.Temperature
	}
	if req.TopP != 0 {
		model.TopP = &req.TopP
	}
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
	}
}

// 设置 genai 历史消息
func setGenAiChatHistoryByOpenaiRequest(cs *genai.ChatSession, req *openModel.ChatCompletionRequest) {
	cs.History = make([]*genai.Content, 0, len(req.Messages))
	if len(req.Messages) > 1 {
		for _, message := range req.Messages[:len(req.Messages)-1] {
			switch message.Role {
			case openai.ChatMessageRoleSystem:
				cs.History = append(cs.History, []*genai.Content{
					{
						Parts: []genai.Part{
							genai.Text(message.StringContent()),
						},
						Role: genModel.GenAiRoleUser,
					},
					{
						Parts: []genai.Part{
							genai.Text("ok."),
						},
						Role: genModel.GenAiRoleModel,
					},
				}...)
			case openai.ChatMessageRoleAssistant:
				cs.History = append(cs.History, &genai.Content{
					Parts: []genai.Part{
						genai.Text(message.StringContent()),
					},
					Role: genModel.GenAiRoleModel,
				})
			case openai.ChatMessageRoleUser:
				cs.History = append(cs.History, &genai.Content{
					Parts: []genai.Part{
						genai.Text(message.StringContent()),
					},
					Role: genModel.GenAiRoleUser,
				})
			}
		}
	}

	if len(cs.History) != 0 && cs.History[len(cs.History)-1].Role != genModel.GenAiRoleModel {
		cs.History = append(cs.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text("ok."),
			},
			Role: genModel.GenAiRoleModel,
		})
	}
}

// genai 转 openai 流式输出响应
func genAiResponseToStreamCompletionResponse(genAiResp *genai.GenerateContentResponse, respID string, created int64) *openModel.CompletionResponse {
	resp := openModel.CompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-%s", respID),
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   genModel.GeminiPro,
		Choices: make([]openModel.CompletionChoice, 0, len(genAiResp.Candidates)),
	}

	for i, candidate := range genAiResp.Candidates {
		var content string
		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			if s, ok := candidate.Content.Parts[0].(genai.Text); ok {
				content = string(s)
			}
		}

		choice := openModel.CompletionChoice{
			Index: i,
		}
		choice.Delta.Content = content

		if candidate.FinishReason > genai.FinishReasonStop {
			log.Printf("genai message finish reason %s\n", candidate.FinishReason.String())

			var openaiFinishReason string = string(openai.FinishReasonStop)
			if candidate.FinishReason == genai.FinishReasonMaxTokens {
				openaiFinishReason = string(openai.FinishReasonLength)
			}
			choice.FinishReason = &openaiFinishReason
		}

		resp.Choices = append(resp.Choices, choice)
	}
	return &resp
}
