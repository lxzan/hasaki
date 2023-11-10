package pb

import (
	"github.com/lxzan/hasaki"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEncoder_ContentType(t *testing.T) {
	assert.Equal(t, Encoder.ContentType(), hasaki.MimeYaml)
}

func TestEncoder_Encode(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var params = struct {
			Name string
			Age  int
		}{Name: "aha"}

		//		var text = `
		//user:
		//  name: "caster"
		//  age: 1
		//`
		_, err := Encoder.Encode(params)
		assert.NoError(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		_, err := Encoder.Encode(nil)
		assert.NoError(t, err)
	})
}

func TestDecode(t *testing.T) {
	var params = struct {
		User struct {
			Name string
			Age  int
		}
	}{}

	var text = `
user:
  name: "caster"
  age: 1
`
	var err = Decode(strings.NewReader(text), &params)
	assert.NoError(t, err)
	assert.Equal(t, params.User.Name, "caster")
}
