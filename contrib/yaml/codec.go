package yaml

import (
	"bytes"
	"github.com/lxzan/hasaki"
	"github.com/lxzan/hasaki/internal"
	"github.com/pkg/errors"
	"github.com/valyala/bytebufferpool"
	"gopkg.in/yaml.v3"
	"io"
)

//go:generate go mod tidy

var Encoder = new(encoder)

type encoder struct{}

func (c encoder) Encode(v any) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}
	w := bytebufferpool.Get()
	err := yaml.NewEncoder(w).Encode(v)
	r := &internal.CloserWrapper{B: w, R: bytes.NewReader(w.B)}
	return r, errors.WithStack(err)
}

func (c encoder) ContentType() string {
	return hasaki.MimeYaml
}

func Decode(r io.Reader, v any) error {
	return errors.WithStack(yaml.NewDecoder(r).Decode(v))
}
