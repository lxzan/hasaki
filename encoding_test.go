package hasaki

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForm_encoder_Encode(t *testing.T) {
	_, err := FormEncoder.Encode(struct {
		Name string
	}{Name: "caster"})
	assert.NoError(t, err)
}
