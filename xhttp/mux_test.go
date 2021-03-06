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
	Pattern  string    `json:"pattern,omitempty"`
	Path     string    `json:"path,omitempty"`
	Error    string    `json:"error,omitempty"`
	Code     string    `json:"code,omitempty"`
	Captures *Captures `json:"captures,omitempty"`
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

type captureHandler struct{}

func (captureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pattern, ok := ctx.Value(RegexpMatchContextKey).(string)
	if !ok {
		pattern = ""
	}

	captures, ok := ctx.Value(CapturesContextKey).(*Captures)
	if !ok {
		InternalError(w, r)
		return
	}

	var data = responseData{
		Pattern:  pattern,
		Path:     requestPathCleanUp(r.URL.Path),
		Captures: captures,
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

	t.Run("pattern with captures", func(t *testing.T) {
		expRegex := []string{
			`^/blog/(?P<blogID>\d{1,19})`,
			`^/blog/(?P<blogUID>[\w-_]+)/images/(?P<imageID>\d{1,19})/thumbnail`,
			`^/blog/(?P<blogUID>[\w-_]+)/images`,
			`^/blog/(?P<blogUID>[\w-_]+)`,
		}
		mux := ServeReMux{}
		mux.Handle(`^/blog/<int:blogID>`, captureHandler{})
		mux.Handle(`^/blog/<blogUID>/images/<int:imageID>/thumbnail`, captureHandler{})
		mux.Handle(`^/blog/<blogUID>/images`, captureHandler{})
		mux.Handle(`^/blog/<str:blogUID>`, captureHandler{})

		for i, h := range mux.handlers.Values() {
			xt.Eq(t, expRegex[i], h.(*reHandler).regex)
		}

		cases := map[string]struct {
			path string
			exp  *Captures
		}{
			"blog 1234 explicit converter": {path: "/blog/1234", exp: &Captures{
				"blogID": Capture{Name: "blogID", Value: "1234", Converter: "int"},
			}},

			"blog Y4sSn8f explicit converter": {path: "/blog/Y4sSn8f", exp: &Captures{
				"blogUID": Capture{Name: "blogUID", Value: "Y4sSn8f", Converter: "str"},
			}},

			"blog Aw5BDfx without converter": {path: "/blog/Aw5BDfx", exp: &Captures{
				"blogUID": Capture{Name: "blogUID", Value: "Aw5BDfx", Converter: "str"},
			}},

			"multiple captures": {path: "/blog/GjqzmwZ/images/2334/thumbnail", exp: &Captures{
				"blogUID": Capture{Name: "blogUID", Value: "GjqzmwZ", Converter: "str"},
				"imageID": Capture{Name: "imageID", Value: "2334", Converter: "int"},
			}},
		}

		for name, cs := range cases {
			t.Run(name, func(t *testing.T) {
				req, err := http.NewRequest(http.MethodGet, cs.path, nil)
				xt.OK(t, err)
				req.Header.Set("Content-Type", ContentTypeJSON)

				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				xt.Eq(t, http.StatusOK, rr.Code)

				var data responseData
				err = json.Unmarshal(rr.Body.Bytes(), &data)
				xt.OK(t, err)
				xt.Eq(t, cs.exp, data.Captures)
			})
		}
	})

	t.Run("named capture cannot occur twice", func(t *testing.T) {
		xt.Panics(t, func() {
			(&reHandler{}).setPattern("/<name1>/foo/<name2>/<name1>")
		})
	})

	t.Run("angle brackets not used correctly when capturing", func(t *testing.T) {
		xt.Panics(t, func() {
			(&reHandler{}).setPattern("/<name1>/foo/<name2")
		})
	})
}
