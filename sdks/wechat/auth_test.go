package wechat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSDK_Code2Session(t *testing.T) {
	sdk := NewSDK("wx6e47a9b5e7a779d0", "3ac6a621a75cb8a88d3bb4d8059583fc")

	resp, err := sdk.Code2Session(&Code2SessionReq{
		JsCode: "123",
	})
	if err != nil {
		t.Error("Code2Session Failed", err)
		t.FailNow()
	}
	t.Logf("resp:%+v\n", resp)
	if !assert.Equal(t, 0, resp.ErrCode) {
		t.FailNow()
	}
}
