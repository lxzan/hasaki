package hasaki

import "testing"

func TestClient_Json(t *testing.T) {
	resp, _ := Get("https://api.github.com/").
		SetProxy("http://127.0.0.1:8888").
		GetBody()
	println(string(resp))
}
