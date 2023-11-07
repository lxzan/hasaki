package hasaki

import (
	"encoding/xml"
	"io"

	jsoniter "github.com/json-iterator/go"
)

func JSONDecode(r io.Reader, v any) error {
	return jsoniter.ConfigFastest.NewDecoder(r).Decode(v)
}

func XMLDecode(r io.Reader, v any) error {
	return xml.NewDecoder(r).Decode(v)
}
