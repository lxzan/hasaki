package hasaki

import "testing"

func TestClient_Json(t *testing.T) {
	resp,_ := Get("https://api.github.com/").
		GetBody()
	println(&resp)
}
