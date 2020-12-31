package wechat

type Response struct {
	ErrCode int64  `json:"errcode" db:"errcode"`
	ErrMsg  string `json:"errmsg" db:"errmsg"`
}
