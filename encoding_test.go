package hasaki

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormEncoder_Encode(t *testing.T) {
	_, err1 := FORMEncoder.Encode(struct {
		Name string
	}{Name: "caster"})
	assert.NoError(t, err1)

	_, err2 := FORMEncoder.Encode(url.Values{
		"name": []string{"caster"},
	})
	assert.NoError(t, err2)
}

func TestJsonEncoder_Encode(t *testing.T) {
	_, err1 := JSONEncoder.Encode(struct {
		Name string
	}{Name: "caster"})
	assert.NoError(t, err1)

	_, err2 := JSONEncoder.Encode(map[string]string{
		"name": "caster",
	})
	assert.NoError(t, err2)
}

func TestStreamEncoder(t *testing.T) {
	encoder := NewStreamEncoder("text/plain")
	assert.Equal(t, encoder.ContentType(), "text/plain")

	_, err1 := encoder.Encode("aha")
	assert.NoError(t, err1)

	_, err2 := encoder.Encode([]byte("aha"))
	assert.NoError(t, err2)

	_, err3 := encoder.Encode(bytes.NewBufferString("oh"))
	assert.NoError(t, err3)

	_, err4 := encoder.Encode(123)
	assert.Error(t, err4)
}
