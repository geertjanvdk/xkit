// Copyright (c) 2020, Geert JM Vanderkelen

package xgraphql

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/eventeneer/xkit/xhttp"
	"github.com/eventeneer/xkit/xt"
)

// MustMD5SumIt returns the hexadecimal string representation of the MD5sum of s.
func MustMD5SumIt(s string) string {
	h := md5.New()
	if _, err := io.WriteString(h, s); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

var responses = map[string]*Response{
	"9714361de333a06aa7ecfac29993eb8e": {
		Data: json.RawMessage(`
{
	"allFilms": {
		"films": [
			{
			  "title": "A New Hope"
			},
			{
			  "title": "The Empire Strikes Back"
			},
			{
			  "title": "Return of the Jedi"
			},
			{
			  "title": "The Phantom Menace"
			},
			{
			  "title": "Attack of the Clones"
			},
			{
			  "title": "Revenge of the Sith"
			},
			{
			  "title": "The Force Awakens"
			},
			{
			  "title": "The Last Jedi"
			},
			{
			  "title": "The Rise of Skywalker"
			}
		]
	}
}
`),
	},
}

func TestNewClient(t *testing.T) {
	t.Run("content type is application/json", func(t *testing.T) {
		c := NewClient("http://example.com/graphql") // URI does not matter for test
		xt.Eq(t, xhttp.ContentTypeJSON, c.ContentType())
	})
}

func TestClient_Execute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlerGraphQLSWAPI))
	defer func() {
		server.Close()
	}()

	t.Run("allFilms", func(t *testing.T) {
		q := `query {
  allFilms {
    films {
      title
    }
  }
}`
		c := NewClient(server.URL)

		var data struct {
			AllFilms struct {
				Films []struct {
					Title string `json:"title"`
				} `json:"films"`
			} `json:"allFilms"`
		}

		err := c.Execute(q, &data)
		xt.OK(t, err)

		var titles []string
		for _, f := range data.AllFilms.Films {
			titles = append(titles, f.Title)
		}
		sort.Strings(titles)
		exp := []string{
			"A New Hope", "Attack of the Clones", "Return of the Jedi",
			"Revenge of the Sith", "The Empire Strikes Back", "The Force Awakens",
			"The Last Jedi", "The Phantom Menace", "The Rise of Skywalker"}
		xt.Eq(t, exp, titles)
	})

	t.Run("error returned", func(t *testing.T) {
		c := NewClient(server.URL)

		var data map[string]interface{}

		err := c.Execute(`{ notExistingQuery { id } }`, &data)
		xt.KO(t, err)
		xt.Eq(t, "test response not available", err.Error())
	})
}

func handlerGraphQLSWAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(xhttp.ContentTypeJSON, r.Header.Get(xhttp.HeaderContentType))
	w.WriteHeader(http.StatusOK)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err.Error()
		_, _ = w.Write([]byte("failed reading body"))
	}

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		_, _ = w.Write([]byte("failed writing body"))
	}

	var response *Response
	response, ok := responses[MustMD5SumIt(payload.Query)]
	if !ok {
		response = &Response{}
		response.Errors = []Error{
			{
				Message: "test response not available",
			},
		}
	}

	buf, err := json.Marshal(response)
	if err != nil {
		_, _ = w.Write([]byte("failed writing body"))
	}

	_, err = w.Write(buf)
	if err != nil {
		_, _ = w.Write([]byte("failed writing body"))
	}
}
