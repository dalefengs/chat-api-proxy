package gemini

import (
	"github.com/dalefeng/chat-api-reverse/global"
	"github.com/dalefeng/chat-api-reverse/model/common/response"
	geminiModel "github.com/dalefeng/chat-api-reverse/model/gemini"
	"github.com/dalefeng/chat-api-reverse/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

type GeMiniApi struct {
}

func (g *GeMiniApi) CompletionsHandler(c *gin.Context) {
	token, err := utils.GetAuthToken(c, "Bearer")
	if err != nil {
		global.SugarLog.Errorw("CompletionsHandler Get hearder token error", "error", err)
		response.FailWithOpenAIError(http.StatusUnauthorized, err.Error(), c)
		return
	}
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(c, option.WithAPIKey(token))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// For text-and-image input (multimodal), use the gemini-pro-vision model
	generativeModel := client.GenerativeModel(geminiModel.GEMINI_PRO_MODEL)
	// Initialize the chat
	cs := generativeModel.StartChat()
	cs.History = []*genai.Content{
		&genai.Content{
			Parts: []genai.Part{
				genai.Text("Hello, I have 2 dogs in my house."),
			},
			Role: "user",
		},
		&genai.Content{
			Parts: []genai.Part{
				genai.Text("Great to meet you. What would you like to know?"),
			},
			Role: "model",
		},
	}

	utils.SetEventStreamHeaders(c)

	iter := cs.SendMessageStream(c, genai.Text("给我上海3日游攻略"))

	for {
		resp, iterErr := iter.Next()
		if iterErr == iterator.Done {
			break
		}
		if iterErr != nil {
			global.SugarLog.Errorw("SendStream iter Error", "error", iterErr, "models", geminiModel.GETMINI_STREAM_METHOD)
		}
		global.SugarLog.Debugw("SendStream iter Next", "Candidates", resp.Candidates, "PromptFeedback", resp.PromptFeedback)
		// ... print resp
	}

}
