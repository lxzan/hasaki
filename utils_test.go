package hasaki

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStructToAny(t *testing.T) {
	as := assert.New(t)
	type Base struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	}

	type base struct {
		Start int64 `json:"start"`
		End   int64 `json:"end"`
	}

	type Model1 struct {
		Base
		Name string `json:"name"`
	}

	type Model2 struct {
		*Base
		Name string `json:"name"`
	}

	type Model3 struct {
		base
		Age  *int   `json:"age"`
		Name string `json:"name"`
	}

	t.Run("struct", func(t *testing.T) {
		var s = Model1{
			Base: Base{Start: 1, End: 2},
			Name: "caster",
		}
		m, err := structToAny(s, "json")
		as.NoError(err)
		as.ElementsMatch([]string{"start", "end", "name"}, getKeys(m))
	})

	t.Run("ptr", func(t *testing.T) {
		var s = Model2{
			Base: &Base{Start: 1, End: 2},
			Name: "caster",
		}
		m, err := structToAny(s, "json")
		as.NoError(err)
		as.ElementsMatch([]string{"start", "end", "name"}, getKeys(m))
	})

	t.Run("private", func(t *testing.T) {
		var age = 12
		var s = Model3{
			base: base{Start: 1, End: 2},
			Age:  &age,
			Name: "caster",
		}
		m, err := structToAny(s, "json")
		as.NoError(err)
		as.ElementsMatch([]string{"name", "age"}, getKeys(m))
	})
}
