// Copyright (c) 2021, Geert JM Vanderkelen

package xhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/geertjanvdk/xkit/xutil"
)

func writeError(w http.ResponseWriter, contentType string, code int, msg string) {
	contentType = strings.TrimSpace(contentType)
	if contentType == "" {
		contentType = ContentTypePlain
	}
	if !strings.Contains(contentType, "charset") &&
		!xutil.HasString([]string{"application/json", "text/plain", "text/html"}, contentType) {
		contentType += "; charset=utf-8"
	}

	var payload string

	switch contentType {
	case ContentTypeJSON:
		var doc = struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		}{
			Error: msg,
			Code:  strconv.FormatInt(http.StatusInternalServerError, 10),
		}
		data, _ := json.Marshal(doc)
		payload = string(data)
	case ContentTypeHTML:
		payload = `<html><body><h3>` + msg + `</h3></body></html>`
	default:
		payload = msg
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, _ = fmt.Fprintln(w, payload)
}

// InternalError replies to the request with a HTTP 500.
func InternalError(w http.ResponseWriter, r *http.Request) {
	msg := "500 internal server error"
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	writeError(w, contentType, http.StatusInternalServerError, msg)
}

// MethodNotAllowed replies to the request with a HTTP 405.
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	msg := "405 method not allowed"
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	writeError(w, contentType, http.StatusMethodNotAllowed, msg)
}

// MethodNotAllowedHandler returns a request handler that replies to each
// request with a ``405 method not allowed''. Depending
// on the request's content type, it will return either
// JSON, HTML or default plain text.
func MethodNotAllowedHandler() http.Handler { return http.HandlerFunc(MethodNotAllowed) }

func NotFound(w http.ResponseWriter, r *http.Request) {
	msg := "404 page not found"
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	writeError(w, contentType, http.StatusNotFound, msg)
}

// NotFoundHandler returns a request handler that replies to each
// request with a ``404 page not found''. Depending
// on the request's content type, it will return either
// JSON, HTML or default plain text.
func NotFoundHandler() http.Handler { return http.HandlerFunc(NotFound) }
