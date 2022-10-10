package hasaki

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormEncode(t *testing.T) {
	as := assert.New(t)
	t.Run("any", func(t *testing.T) {
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
	})

	t.Run("struct", func(t *testing.T) {
		var model = struct {
			Name    string `form:"name"`
			Age     int    `form:"age"`
			Options []int  `form:"options"`
		}{
			Name:    "caster",
			Age:     12,
			Options: []int{1, 2, 3},
		}
		result, err := FormEncoder.Encode(&model)
		as.NoError(err)
		as.Equal("age=12&name=caster&options=1&options=2&options=3", string(result))
	})
}
