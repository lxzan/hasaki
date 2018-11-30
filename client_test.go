package hasaki

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestClient_Json(t *testing.T) {
	data, _ := ioutil.ReadFile("/Users/caster/go/src/github.com/lxzan/hasaki/LICENSE")
	resp, err := Put("http://devfeng-bbs-att-1255531212.cos.ap-shanghai.myqcloud.com/chat/2018/11/30/1444229kk8j6622y2gsfvo.png", nil).
		Set(Form{
			"Authorization": "q-sign-algorithm=sha1&q-ak=AKIDroUF0K02NSaIe3jpCa3fKzeJzeKbulcT&q-sign-time=1543560262;1543560862&q-key-time=1543560262;1543560862&q-header-list=content-length&q-url-param-list=&q-signature=9ac4cbae1c0e8d5461eafe6f1614d9ea3582ef2d",
			"Host":          "devfeng-bbs-att-1255531212.cos.ap-shanghai.myqcloud.com",
		}).
		SetBody(bytes.NewReader(data)).
		GetBody()
	println(string(resp), &err)
}
