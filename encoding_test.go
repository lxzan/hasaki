package hasaki

import "testing"

func TestFormEncode(t *testing.T) {
	if result, err := FormEncode(Any{
		"a": 1,
		"b": "xxx",
		"c": []string{"x", "y", "z"},
		"d": []float64{1.1, 1.21},
	}); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(result)
	}

	if _, err := FormEncode(Any{
		"a": 1,
		"b": "xxx",
		"c": []byte{'x'},
		"d": []float64{1.1, 1.21},
	}); err != nil {
		t.Log(err.Error())
	} else {
		t.Error(err.Error())
	}
}
