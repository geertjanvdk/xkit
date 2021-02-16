// Copyright (c) 2020, Geert JM Vanderkelen

package blackbox

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/geertjanvdk/xkit/xhttp"
	"github.com/geertjanvdk/xkit/xnet"
	"github.com/geertjanvdk/xkit/xt"
)

func runTLSServer(addr net.TCPAddr) {
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("pong"))
	})

	err := http.ListenAndServeTLS(addr.String(),
		"_data/test_server.crt.pem",
		"_data/test_server.key.pem",
		nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}

func TestXHTTPNewClient(t *testing.T) {
	tlsAddr := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: xnet.MustGetLocalhostTCPPort(),
	}

	go runTLSServer(tlsAddr)
	time.Sleep(time.Second)

	t.Run("set and get authZ bearer", func(t *testing.T) {
		b := "my_bearer"
		c := xhttp.NewClient("http://example.com")
		xt.Eq(t, "", c.AuthzBearer())
		c.SetAuthzBearer(b)
		xt.Eq(t, b, c.AuthzBearer())
	})

	t.Run("functional option WithBearer", func(t *testing.T) {
		b := "my_bearer"
		c := xhttp.NewClient("http://example.com", xhttp.WithBearer(b))
		xt.Eq(t, b, c.AuthzBearer())
	})

	t.Run("using TLS without functional option WithTLSInsecure", func(t *testing.T) {
		c := xhttp.NewClient("https://" + tlsAddr.String() + "/ping")
		_, err := c.Get()
		xt.KO(t, err)
		xt.Match(t, ".*x509: cannot validate certificate for .*", err.Error())
	})

	t.Run("using TLS with functional option WithTLSInsecure", func(t *testing.T) {
		c := xhttp.NewClient("https://"+tlsAddr.String()+"/ping", xhttp.WithTLSInsecure())
		req, err := c.Get()
		xt.OK(t, err)
		body, err := ioutil.ReadAll(req.Body)
		xt.OK(t, err)
		xt.Eq(t, "pong", string(body))
	})
}
