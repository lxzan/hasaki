package hasaki

import (
	"bytes"
	"encoding/xml"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestForm_encoder_Encode(t *testing.T) {
	_, err1 := FormCodec.Encode(struct {
		Name string
	}{Name: "caster"})
	assert.True(t, errors.Is(err1, errUnsupportedData))

	_, err2 := FormCodec.Encode(url.Values{
		"name": []string{"caster"},
	})
	assert.NoError(t, err2)

	_, err3 := FormCodec.Encode(nil)
	assert.NoError(t, err3)

	var netConn *net.TCPConn
	_, err4 := FormCodec.Encode(net.Conn(netConn))
	assert.Error(t, err4)

	_, err5 := FormCodec.Encode("a=xxx")
	assert.NoError(t, err5)
}

func TestJson_encoder_Encode(t *testing.T) {
	_, err1 := JsonCodec.Encode(struct {
		Name string
	}{Name: "caster"})
	assert.NoError(t, err1)

	_, err2 := JsonCodec.Encode(map[string]interface{}{
		"name": "caster",
	})
	assert.NoError(t, err2)

	_, err3 := JsonCodec.Encode(nil)
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
		var err = FormCodec.Decode(strings.NewReader(text), &params)
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
		var err = FormCodec.Decode(strings.NewReader(""), params)
		assert.True(t, errors.Is(err, errUnsupportedData))
	})

	t.Run("error", func(t *testing.T) {
		var params = url.Values{}
		var text = "a;b;c"
		var err = FormCodec.Decode(strings.NewReader(text), &params)
		assert.Error(t, err)
	})
}

func TestFormDecodeReadError(t *testing.T) {
	readErr := errors.New("form read failed")
	for _, tc := range []struct {
		name   string
		reader io.Reader
	}{
		{"empty", iotest.ErrReader(readErr)},
		{"partial form", io.MultiReader(strings.NewReader("a=1"), iotest.ErrReader(readErr))},
		{"malformed form", io.MultiReader(strings.NewReader("a=%"), iotest.ErrReader(readErr))},
		{"data and error", &formReadErrorReader{data: "a=1", err: readErr}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			params := url.Values{"existing": {"value"}}
			err := FormCodec.Decode(tc.reader, &params)
			assert.ErrorIs(t, err, readErr)
			assert.Equal(t, url.Values{"existing": {"value"}}, params)
		})
	}
}

func TestFormDecodeEOF(t *testing.T) {
	for _, tc := range []struct {
		name string
		data string
		want url.Values
	}{
		{"empty", "", url.Values{}},
		{"form", "a=1&a=2&b=hello+world", url.Values{"a": {"1", "2"}, "b": {"hello world"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			params := url.Values{"existing": {"value"}}
			err := FormCodec.Decode(&formReadErrorReader{data: tc.data, err: io.EOF}, &params)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, params)
		})
	}
}

func TestBindFormTruncatedResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", MimeForm)
		w.Header().Set("Content-Length", "10")
		_, _ = io.WriteString(w, "a=1")
	}))
	defer srv.Close()
	client, err := NewClient(WithHTTPClient(srv.Client()))
	assert.NoError(t, err)
	resp := client.Get(srv.URL).Send(nil)
	if !assert.NoError(t, resp.Err()) {
		return
	}
	defer resp.Body.Close()
	params := url.Values{"existing": {"value"}}
	err = resp.BindForm(&params)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
	assert.Equal(t, url.Values{"existing": {"value"}}, params)
}

type formReadErrorReader struct {
	data string
	err  error
}

func (r *formReadErrorReader) Read(p []byte) (int, error) {
	n := copy(p, r.data)
	r.data = r.data[n:]
	if len(r.data) == 0 {
		return n, r.err
	}
	return n, nil
}

func TestXmlEncoder(t *testing.T) {
	t.Run("type", func(t *testing.T) {
		assert.Equal(t, XmlCodec.ContentType(), MimeXml)
	})

	t.Run("", func(t *testing.T) {
		var params = struct {
			XMLName xml.Name `xml:"xml"`
			Name    string   `xml:"name"`
			Age     int      `xml:"age"`
		}{
			Name: "cas",
		}
		_, err := XmlCodec.Encode(params)
		assert.NoError(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		_, err := XmlCodec.Encode(nil)
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
	var err = XmlCodec.Decode(strings.NewReader(text), p)
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
