// Copyright (c) 2021, Geert JM Vanderkelen

package xhttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geertjanvdk/xkit/xt"
)

type pathEchoHandler struct{}

type responseData struct {
	Pattern string `json:"pattern,omitempty"`
	Path    string `json:"path,omitempty"`
	Error   string `json:"error,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (pathEchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pattern, ok := ctx.Value(RegexpMatchContextKey).(string)
	if !ok {
		pattern = ""
	}

	var data = responseData{
		Pattern: pattern,
		Path:    requestPathCleanUp(r.URL.Path),
	}
	body, err := json.Marshal(data)
	if err != nil {
		InternalError(w, r)
		return
	}
	_, _ = w.Write(body)
}

func TestServeReMux_Handle(t *testing.T) {
	t.Run("must compile regex", func(t *testing.T) {
		mux := ServeReMux{}
		mux.Handle("^/?$", nil)
		xt.Eq(t, 1, mux.handlers.Count())
	})

	t.Run("compile of regex panics", func(t *testing.T) {
		xt.Panics(t, func() {
			NewServeReMux().Handle(`\C`, nil) // Go does not like \C
		})
	})

	t.Run("route answers on foo and foobar", func(t *testing.T) {
		mux := ServeReMux{}
		mux.Handle(`^/fo(o|bar)$`, pathEchoHandler{})
		xt.Eq(t, 1, mux.handlers.Count())

		cases := map[string]struct {
			expCode int
		}{
			"/foo":     {expCode: http.StatusOK},
			"/fobar":   {expCode: http.StatusOK},
			"/fo":      {expCode: http.StatusNotFound},
			"/foobarr": {expCode: http.StatusNotFound},
			"/fooba":   {expCode: http.StatusNotFound},
		}

		for p, cs := range cases {
			t.Run(p, func(t *testing.T) {
				req, err := http.NewRequest(http.MethodGet, p, nil)
				xt.OK(t, err)
				req.Header.Set("Content-Type", ContentTypeJSON)

				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				var data responseData
				err = json.Unmarshal(rr.Body.Bytes(), &data)
				xt.OK(t, err)

				xt.Eq(t, cs.expCode, rr.Code)
			})
		}
	})

	t.Run("nested routes", func(t *testing.T) {
		mux := ServeReMux{}
		mux.Handle(`/fo(o|bar)/`, pathEchoHandler{})
		xt.Eq(t, 1, mux.handlers.Count())

		cases := map[string]struct {
			expCode int
		}{
			"/foo/a/b/c":     {expCode: http.StatusOK},
			"/a/b/c/fobar/d": {expCode: http.StatusOK},
			"/fo/bar":        {expCode: http.StatusNotFound},
			"/foobarr":       {expCode: http.StatusNotFound},
			"/fooba":         {expCode: http.StatusNotFound},
		}

		for p, cs := range cases {
			t.Run(p, func(t *testing.T) {
				req, err := http.NewRequest(http.MethodGet, p, nil)
				xt.OK(t, err)
				req.Header.Set("Content-Type", ContentTypeJSON)

				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				var data responseData
				err = json.Unmarshal(rr.Body.Bytes(), &data)
				xt.OK(t, err)

				xt.Eq(t, cs.expCode, rr.Code)
			})
		}
	})

	t.Run("first come first served", func(t *testing.T) {
		mux := ServeReMux{}
		mux.Handle(`^/foo/bar`, pathEchoHandler{})
		mux.Handle(`foo`, pathEchoHandler{})
		mux.Handle(`^/$`, pathEchoHandler{})
		xt.Eq(t, 3, mux.handlers.Count())

		cases := map[string]struct {
			expCode int
			exp     string
		}{
			"/foo":            {expCode: http.StatusOK, exp: `foo`},
			"/a/b/c/fobar/d":  {expCode: http.StatusNotFound},
			"/a/b/c/foobar/d": {expCode: http.StatusOK, exp: `foo`},
			"foo/bar":         {expCode: http.StatusOK, exp: `foo`},
			"/foo/bar":        {expCode: http.StatusOK, exp: `^/foo/bar`},
			"/barr":           {expCode: http.StatusNotFound},
			"/fo":             {expCode: http.StatusNotFound},
			"/":               {expCode: http.StatusOK, exp: `^/$`},
		}

		for p, cs := range cases {
			t.Run(p+" path", func(t *testing.T) {
				req, err := http.NewRequest(http.MethodGet, p, nil)
				xt.OK(t, err)
				req.Header.Set("Content-Type", ContentTypeJSON)

				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				var data responseData
				err = json.Unmarshal(rr.Body.Bytes(), &data)
				xt.OK(t, err)

				xt.Eq(t, cs.expCode, rr.Code, "matched pattern: "+data.Pattern)
			})
		}
	})

	t.Run("specifying no methods only allows GET", func(t *testing.T) {
		mux := ServeReMux{}
		mux.Handle(`^/nomethod-defaultget`, pathEchoHandler{})

		req, err := http.NewRequest(http.MethodPost, "/nomethod-defaultget",
			strings.NewReader("return to sender"))
		xt.OK(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		xt.Eq(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("specifying no methods only allows GET", func(t *testing.T) {
		mux := ServeReMux{}
		mux.Handle(`^/nomethod-defaultget`, pathEchoHandler{})

		req, err := http.NewRequest(http.MethodPost, "/nomethod-defaultget",
			strings.NewReader("return to sender"))
		xt.OK(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		xt.Eq(t, http.StatusMethodNotAllowed, rr.Code)
	})

	t.Run("a POST only route", func(t *testing.T) {
		mux := ServeReMux{}
		mux.Handle(`^/postonly`, pathEchoHandler{}, MethodPost)

		req, err := http.NewRequest(http.MethodGet, "/postonly", nil)
		xt.OK(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		xt.Eq(t, http.StatusMethodNotAllowed, rr.Code)
	})
}
