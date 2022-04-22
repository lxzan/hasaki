package hasaki

import "testing"

func TestFormEncode(t *testing.T) {
	if result, err := FormEncoder.Encode(Any{
		"a": 1,
		"b": "xxx",
		"c": []string{"x", "y", "z"},
		"d": []float64{1.1, 1.21},
	}); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(string(result))
	}

	if result, err := FormEncoder.Encode(Any{
		"a": 1,
		"b": "xxx",
		"c": []byte{'x'},
		"d": []float64{1.1, 1.21},
	}); err == nil {
		t.Log(string(result))
	} else {
		t.Error(err.Error())
	}
}
