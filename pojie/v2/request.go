package v2

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// RequestPayload 定义了整个批量请求的结构体
type RequestPayload struct {
	Requests []SubRequest `json:"requests"`
}

// SubRequest 定义了单个子请求的结构体
type SubRequest struct {
	Method  string            `json:"method"`
	Timeout int               `json:"timeout"`
	Header  map[string]string `json:"header"`
	Path    string            `json:"path"`
	Body    string            `json:"body"`
}

// BuildBatchRequestPayload 构造并返回批量请求的 JSON 字节数组
// 现在它接收一个 password 参数
func BuildBatchRequestPayload(qrcodeRouteParam, dynamicPassword string) ([]byte, error) {
	// 对密码进行URL编码，以确保在 form-urlencoded 格式中正确传递
	encodedPassword := url.QueryEscape(dynamicPassword)

	// 定义通用的 qrcode_route，因为它在所有请求中都是相同的
	qrcodeRoute := url.QueryEscape(qrcodeRouteParam) // 同样需要URL编码

	// 构建第一个请求的 body 字符串，包含动态密码
	// 注意：其他固定参数如 render_default_fields 等也需要URL编码
	firstRequestBody := fmt.Sprintf(
		"qrcode_route=%s&password=%s&render_default_fields=0&render_component_number=0&render_edit_btn=1&package_id=",
		qrcodeRoute,
		encodedPassword,
	)

	// 构建第二个请求的 body 字符串，如果 password=null 是固定值，则直接用字符串 "null"
	// 如果 password 也可以是动态的 "null" 字符串，则可以重复使用 encodedPassword
	secondRequestBody := fmt.Sprintf(
		"qrcode_route=%s&password=null&env=h5&org_id=undefined&package_id=",
		qrcodeRoute,
	)

	// 其他请求的 body 保持不变，但使用编码后的 qrcodeRoute
	thirdRequestBody := fmt.Sprintf(
		"qrcode_route=%s&with_state_change_log=1&package_id=",
		qrcodeRoute,
	)
	fourthRequestBody := fmt.Sprintf(
		"qrcode_route=%s&env=h5&scan_source=1&package_id=",
		qrcodeRoute,
	)
	fifthRequestBody := fmt.Sprintf(
		"qrcode_route=%s&package_id=",
		qrcodeRoute,
	)

	payload := RequestPayload{
		Requests: []SubRequest{
			{
				Method:  "POST",
				Timeout: 10000,
				Header:  map[string]string{"content-type": "application/x-www-form-urlencoded"},
				Path:    "/qrcoderoute/qrcodeRouteNew",
				Body:    firstRequestBody, // 使用动态密码构建的 body
			},
			{
				Method:  "POST",
				Timeout: 10000,
				Header:  map[string]string{"content-type": "application/x-www-form-urlencoded"},
				Path:    "/record/getRecordTpl",
				Body:    secondRequestBody, // 这里的 password 仍然是 "null"
			},
			{
				Method:  "POST",
				Timeout: 10000,
				Header:  map[string]string{"content-type": "application/x-www-form-urlencoded"},
				Path:    "/state/getTargetStateMsg",
				Body:    thirdRequestBody,
			},
			{
				Method:  "POST",
				Timeout: 10000,
				Header:  map[string]string{"content-type": "application/x-www-form-urlencoded"},
				Path:    "/operation/getQrcodeOperationByQrcodeRoute",
				Body:    fourthRequestBody,
			},
			{
				Method:  "POST",
				Timeout: 10000,
				Header:  map[string]string{"content-type": "application/x-www-form-urlencoded"},
				Path:    "/qrcoderoute/getQrcodeProperties",
				Body:    fifthRequestBody,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON payload: %w", err)
	}
	return jsonPayload, nil
}
