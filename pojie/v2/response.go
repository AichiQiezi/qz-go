package v2

import (
	"encoding/json"
	"fmt"
)

// ResponsePayload 定义了整个批量响应的结构体
type ResponsePayload struct {
	Responses []SubResponse `json:"responses"`
}

// SubResponse 定义了单个子响应的结构体
type SubResponse struct {
	Status int    `json:"status"`
	Reason string `json:"reason"`
	Body   string `json:"body"` // 注意：这里 body 是一个 JSON 字符串，你需要再次解析它
	URL    string `json:"url"`
	Wait   int    `json:"wait"`
}

// ParseBatchResponse 解析响应体并返回 ResponsePayload 结构体
func ParseBatchResponse(respBody []byte) (*ResponsePayload, error) {
	var responsePayload ResponsePayload
	err := json.Unmarshal(respBody, &responsePayload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response JSON: %w", err)
	}
	return &responsePayload, nil
}
