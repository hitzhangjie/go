package tests

import (
	"bytes"
	"encoding/json"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_JSON_Marshal(t *testing.T) {
	m := Email{
		From:      "",
		To:        "",
		Title:     "",
		Body:      "",
		Signature: "",
	}
	_, err := json.Marshal(m)
	require.Nil(t, err)
}

type Email struct {
	From      string `json:"email,omitempty"`
	To        string `json:"to,omitempty"`
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type EmailX struct {
	From      string `json:"email,omitempty"`
	To        string `json:"to,omitempty"`
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
	Signature string `json:"signature,omitempty"`
}

var pool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// func (ex *EmailX) MarshalJSON() ([]byte, error) {
// 	v := pool.Get().([]byte)
// 	defer pool.Put(v)

// 	b := bytes.NewBuffer(v)
// 	//fmt.Fprintf(b, `{"from":"%s","to":"%s","title":"%s","body":"%s","signature":"%s"}`, ex.From, ex.To, ex.Title, ex.Body, ex.Signature)
// 	fmt.Fprintf(b, `{"from":"%s"}`, ex.From)
// 	return b.Bytes(), nil
// }

// func (ex *EmailX) MarshalJSON() ([]byte, error) {
// 	buf := pool.Get().(*bytes.Buffer)
// 	buf.Reset()
// 	defer pool.Put(buf)

// 	buf.WriteString(`{`)
// 	buf.WriteString(`"from":"` + ex.From + `",`)
// 	buf.WriteString(`"to":"` + ex.To + `",`)
// 	buf.WriteString(`"title":"` + ex.Title + `",`)
// 	buf.WriteString(`"body":"` + ex.Body + `",`)
// 	buf.WriteString(`"signature":"` + ex.Signature + `"`)
// 	buf.WriteString(`}`)
// 	return buf.Bytes(), nil
// }

func (ex *EmailX) MarshalJSON() ([]byte, error) {

	// Pre-allocate byte slice to avoid extra allocations
	buf := make([]byte, 0, 128)

	buf = append(buf, '{')

	// Use strconv.Quote to get quoted strings instead of concatenating
	buf = append(buf, []byte(`"from":`)...)
	buf = strconv.AppendQuote(buf, ex.From)
	buf = append(buf, ',')

	buf = append(buf, []byte(`"to":`)...)
	buf = strconv.AppendQuote(buf, ex.To)
	buf = append(buf, ',')

	buf = append(buf, []byte(`"title":`)...)
	buf = strconv.AppendQuote(buf, ex.Title)
	buf = append(buf, ',')

	buf = append(buf, []byte(`"body":`)...)
	buf = strconv.AppendQuote(buf, ex.Body)
	buf = append(buf, ',')

	buf = append(buf, []byte(`"signature":`)...)
	buf = strconv.AppendQuote(buf, ex.Signature)

	buf = append(buf, '}')

	return buf, nil
}

func Benchmark_JSON_Marshal(b *testing.B) {
	m := Email{
		From:      "abc@qq.com",
		To:        "def@qq.com",
		Title:     "xyz",
		Body:      "def",
		Signature: "lll",
	}
	b.Run("Email: not implement Marshaler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(&m)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	mm := EmailX{
		From:      "abc@qq.com",
		To:        "def@qq.com",
		Title:     "xyz",
		Body:      "def",
		Signature: "lll",
	}
	b.Run("Email: implement Marshaler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(&mm)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
