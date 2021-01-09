package wechat

import "net/http"

type Code2SessionReq struct {
	JsCode string
}
type Code2SessionRsp struct {
	*Response
	OpenId     string `json:"openid" db:"openid"`
	SessionKey string `json:"session_key" db:"session_key"`
	UnionId    string `json:"unionid" db:"unionid"`
}

func (s *SDK) Code2Session(req *Code2SessionReq) (*Code2SessionRsp, error) {
	params := map[string]string{
		"js_code":    req.JsCode,
		"grant_type": "authorization_code",
		"appid":      s.AppId,
		"secret":     s.AppSecret,
	}
	r, err := http.NewRequest(http.MethodGet, s.getURL("/sns/jscode2session", params), nil)
	if err != nil {
		return nil, err
	}

	resp := new(Code2SessionRsp)
	if err := s.getJson(r, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
