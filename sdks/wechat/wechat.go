package wechat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SDK struct {
	Host      string
	AppId     string
	AppSecret string
}

func NewSDK(appId, appSecret string) *SDK {
	return &SDK{
		Host:      "https://api.weixin.qq.com",
		AppId:     appId,
		AppSecret: appSecret,
	}
}

func (s *SDK) getJson(req *http.Request, ret interface{}) error {
	resp, err := s.get(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, ret)
}

func (s *SDK) get(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func (s *SDK) getURL(path string, params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	return s.Host + path + "?" + values.Encode()
}
