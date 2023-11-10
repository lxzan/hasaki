package hasaki

import (
	"bytes"
	"encoding/xml"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForm_encoder_Encode(t *testing.T) {
	_, err1 := FormEncoder.Encode(struct {
		Name string
	}{Name: "caster"})
	assert.True(t, errors.Is(err1, errUnsupportedData))

	_, err2 := FormEncoder.Encode(url.Values{
		"name": []string{"caster"},
	})
	assert.NoError(t, err2)

	_, err3 := FormEncoder.Encode(nil)
	assert.NoError(t, err3)

	var netConn *net.TCPConn
	_, err4 := FormEncoder.Encode(net.Conn(netConn))
	assert.Error(t, err4)

	_, err5 := FormEncoder.Encode("a=xxx")
	assert.NoError(t, err5)
}

func TestJson_encoder_Encode(t *testing.T) {
	_, err1 := JsonEncoder.Encode(struct {
		Name string
	}{Name: "caster"})
	assert.NoError(t, err1)

	_, err2 := JsonEncoder.Encode(map[string]interface{}{
		"name": "caster",
	})
	assert.NoError(t, err2)

	_, err3 := JsonEncoder.Encode(nil)
	assert.NoError(t, err3)
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

func TestFormDecode(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var params = url.Values{}
		var text = "a=xxx&b=1"
		var err = FormDecode(strings.NewReader(text), &params)
		assert.NoError(t, err)
		assert.Equal(t, params["a"][0], "xxx")
	})

	t.Run("bind", func(t *testing.T) {
		var text = "a=xxx&b=1"
		var resp = &Response{
			Response: &http.Response{
				Body: io.NopCloser(strings.NewReader(text)),
			},
		}
		var params = url.Values{}
		var err = resp.BindForm(&params)
		assert.NoError(t, err)
		assert.Equal(t, params["a"][0], "xxx")
	})

	t.Run("unsupported type", func(t *testing.T) {
		var params = url.Values{}
		var err = FormDecode(strings.NewReader(""), params)
		assert.True(t, errors.Is(err, errUnsupportedData))
	})

	t.Run("error", func(t *testing.T) {
		var params = url.Values{}
		var text = "a;b;c"
		var err = FormDecode(strings.NewReader(text), &params)
		assert.Error(t, err)
	})
}

func TestXmlEncoder(t *testing.T) {
	t.Run("type", func(t *testing.T) {
		assert.Equal(t, XmlEncoder.ContentType(), MimeXml)
	})

	t.Run("", func(t *testing.T) {
		var params = struct {
			XMLName xml.Name `xml:"xml"`
			Name    string   `xml:"name"`
			Age     int      `xml:"age"`
		}{
			Name: "cas",
		}
		_, err := XmlEncoder.Encode(params)
		assert.NoError(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		_, err := XmlEncoder.Encode(nil)
		assert.NoError(t, err)
	})
}

func TestXmlDecode(t *testing.T) {
	var text = `
<?xml version="1.0" encoding="UTF-8" ?>
<peoples version="0.9">
    <people id="888">
        <name>msr</name>
        <address>中国上海</address>
    </people>
    <people id="998">
        <name>maishuren</name>
        <address>中国上海</address>
    </people>
</peoples>
`
	type Peoples struct {
		XMLName xml.Name `xml:"peoples"`
		Version string   `xml:"version,attr"`
		Peos    []struct {
			XMLName xml.Name `xml:"people"`
			Id      int      `xml:"id,attr"`
			Name    string   `xml:"name"`
			Address string   `xml:"address"`
		} `xml:"people"`
	}

	var p = &Peoples{}
	var err = XmlDecode(strings.NewReader(text), p)
	assert.NoError(t, err)

	t.Run("bind", func(t *testing.T) {
		var resp = &Response{
			Response: &http.Response{
				Body: io.NopCloser(strings.NewReader(text)),
			},
		}
		var params = Peoples{}
		var err = resp.BindXML(&params)
		assert.NoError(t, err)
		assert.Equal(t, params.Peos[0].Id, 888)
	})
}
