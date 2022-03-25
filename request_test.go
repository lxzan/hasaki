package hasaki

import "testing"

func TestNewRequest(t *testing.T) {
	const baseurl = "http://localhost:9000"
	if _, err := Get(baseurl + "/p1").Send(nil).GetBody(); err != nil {
		t.Error(err.Error())
	}

	if _, err := Get(baseurl + "/p2").Send(nil).GetBody(); err == nil {
		t.Error(err.Error())
	}

	if _, err := Get(baseurl + "/p3").Send(nil).GetBody(); err == nil {
		t.Error(err.Error())
	}

	if _, err := Get(baseurl + "/p4").Send(nil).GetBody(); err == nil {
		t.Error(err.Error())
	}
}
