package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/linkingthing/ddi-monitor/config"
)

func GetToken(client *http.Client, controllerAddr string) (string, error) {
	loginRequest := struct {
		Username string
		Password string
	}{
		Username: config.Username,
		Password: config.Password,
	}

	return HttpRequest(client, http.MethodPost,
		fmt.Sprintf("https://%s%s", controllerAddr, "/apis/linkingthing.com/common/v1/getsystemapitoken"),
		"", &loginRequest)
}

func GenControllerRequestUrl(controllerAddr, action, id string) string {
	var builder strings.Builder
	builder.WriteString("https://")
	builder.WriteString(controllerAddr)
	builder.WriteString("/apis/linkingthing.com/metric/v1/nodes/")
	builder.WriteString(id)
	builder.WriteString("?action=")
	builder.WriteString(action)
	return builder.String()
}

func HttpRequest(cli *http.Client, httpMethod, url, token string, req interface{}) (string, error) {
	var httpReqBody io.Reader
	if req != nil {
		reqBody, err := json.Marshal(req)
		if err != nil {
			return "", fmt.Errorf("marshal request failed: %s", err.Error())
		}

		httpReqBody = bytes.NewBuffer(reqBody)
	}

	httpReq, err := http.NewRequest(httpMethod, url, httpReqBody)
	if err != nil {
		return "", fmt.Errorf("new http request failed: %s", err.Error())
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if token != "" {
		httpReq.Header.Set(config.AuthKey, token)
	}

	httpResp, err := cli.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("send http request failed: %s", err.Error())
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("read http response body failed: %s", err.Error())
	}

	if httpResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get wrong response with code %v: %s", httpResp.StatusCode, body)
	}

	return httpResp.Header.Get(config.AuthKey), nil
}
