package v2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"qz-go/pojie/db"
)

func TryPassword2(qrcodeRoute, password string) (bool, error) {
	baseUrl := "https://nc.caoliao.net/batch-requests"

	requestPayload, err := BuildBatchRequestPayload(qrcodeRoute, password)
	if err != nil {
		return false, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(requestPayload))
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	rps, err := ParseBatchResponse(body)
	if err != nil {
		return false, err
	}

	rp := rps.Responses[0]
	return checkSuccess([]byte(rp.Body)), nil
}

func checkSuccess(body []byte) bool {
	var res struct {
		Msg struct {
			Text string `json:"text"`
			Code string `json:"code"`
		} `json:"msg"`
	}

	err := json.Unmarshal(body, &res)
	if err != nil {
		db.LogDBError("解析body失败")
		return false
	}

	return res.Msg.Code == "0"
}
