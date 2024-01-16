package token

import "time"

type TokensRequest struct {
	Token     string   `json:"token"`
	Copilot   []string `json:"copilot"`
	CoCopilot []string `json:"cocopilot"`
	ChatGPT   []string `json:"chatgpt"`
	Gemini    []string `json:"gemini"`
}

type PoolTokenInfo struct {
	Token      string                    `json:"token"`       // token
	Models     map[string]*PoolModelInfo `json:"models"`      // 模型详情信息
	CreateTime time.Time                 `json:"create_time"` // 创建时间
	UpdateTime *time.Time                `json:"update_time"` // 更新时间
	LastTime   *time.Time                `json:"last_time"`   // 最后一次调用时间
}

type PoolModelInfo struct {
	Name       string         `json:"name"`
	Count      int            `json:"count"` // token 数量
	Tokens     []TokenInfo    `json:"tokens"`
	Index      int            `json:"index"`       // 当前使用的 token 索引
	Usage      int            `json:"usage"`       // 总调用次数
	DayUsage   int            `json:"day_usage"`   // 当天调用次数
	ModelUsage map[string]int `json:"model_usage"` // 模型 调用次数
	LastTime   *time.Time     `json:"last_time"`   // 最后一次调用时间
}

type TokenInfo struct {
	Token    string     `json:"token"`
	Expire   int        `json:"expire"`    // 过期时间
	LastTime *time.Time `json:"last_time"` // 最后一次调用时间
}

func NewPoolModelInfo(tokens []string) *PoolModelInfo {
	tokensInfos := make([]TokenInfo, len(tokens))
	for i, token := range tokens {
		tokensInfos[i] = TokenInfo{
			Token:    token,
			Expire:   0,
			LastTime: nil,
		}
	}
	return &PoolModelInfo{
		Tokens:     tokensInfos,
		Count:      len(tokensInfos),
		Index:      0,
		Usage:      0,
		DayUsage:   0,
		ModelUsage: make(map[string]int),
		LastTime:   nil,
	}
}
