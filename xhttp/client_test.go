// Copyright (c) 2020, Geert JM Vanderkelen

package xhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

func TestNewClient(t *testing.T) {
	t.Run("default content type", func(t *testing.T) {
		xt.Eq(t, ContentTypeJSON, defaultContentType, "expected correct default content type")
		c := NewClient("http://example.com")
		xt.Eq(t, ContentTypeJSON, c.ContentType())
	})

	t.Run("set different content-type", func(t *testing.T) {
		c := NewClient("http://example.com", WithContentType(ContentTypePlain))
		xt.Eq(t, ContentTypePlain, c.ContentType())
	})

	t.Run("method GET", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(handlerGetHello))
		defer func() {
			server.Close()
		}()

		c := NewClient(server.URL)
		resp, err := c.Get()
		xt.OK(t, err)
		xt.Eq(t, http.StatusOK, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		xt.OK(t, err)
		xt.Eq(t, `"hello!"`, string(body))
		xt.Eq(t, defaultContentType, resp.Header.Get(HeaderContentType))
	})

	t.Run("method POST", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(handlerPostEcho))
		defer func() {
			server.Close()
		}()

		exp := `"echo!"`

		c := NewClient(server.URL)
		resp, err := c.Post(bytes.NewReader([]byte(exp)))
		xt.OK(t, err)
		xt.Eq(t, http.StatusOK, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		xt.OK(t, err)
		xt.Eq(t, exp, string(body))
		xt.Eq(t, defaultContentType, resp.Header.Get(HeaderContentType))
	})

	t.Run("authorization bearer", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(handlerPostAuthorization))
		defer func() {
			server.Close()
		}()

		exp := "this.is.my.authorization.bearer"

		c := NewClient(server.URL, WithBearer(exp))
		resp, err := c.Post(bytes.NewReader([]byte("not used in test")))
		xt.OK(t, err)
		xt.Eq(t, http.StatusOK, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		xt.OK(t, err)
		xt.Eq(t, defaultContentType, resp.Header.Get(HeaderContentType))

		var data struct {
			Bearer string `json:"bearer"`
		}
		xt.OK(t, json.Unmarshal(body, &data))
		xt.Eq(t, "Bearer "+exp, data.Bearer)
	})
}

func TestClient_AuthBearer(t *testing.T) {
	c := NewClient("http://127.0.0.1")
	exp := "this.is.my.authorization.bearer"
	c.SetAuthzBearer(exp)
	xt.Eq(t, exp, c.AuthzBearer())
}

func handlerGetHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderContentType, r.Header.Get(HeaderContentType))
	w.WriteHeader(http.StatusOK)

	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = w.Write([]byte("failed reading body"))
	}

	payload := []byte(`"hello!"`)
	_, err = w.Write(payload)
	if err != nil {
		_, _ = w.Write([]byte("failed writing body"))
	}
}

func handlerPostEcho(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderContentType, r.Header.Get(HeaderContentType))
	w.WriteHeader(http.StatusOK)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = w.Write([]byte("failed reading body"))
	}

	_, err = w.Write(body)
	if err != nil {
		_, _ = w.Write([]byte("failed writing body"))
	}
}

func handlerPostAuthorization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(HeaderContentType, r.Header.Get(HeaderContentType))
	w.WriteHeader(http.StatusOK)

	bearer := r.Header.Get(HeaderAuthorization)

	payload, err := json.Marshal(struct {
		Bearer string `json:"bearer"`
	}{
		Bearer: bearer,
	})
	if err != nil {
		_, _ = w.Write([]byte("failed writing body"))
	}

	_, err = w.Write(payload)
	if err != nil {
		_, _ = w.Write([]byte("failed writing body"))
	}
}
