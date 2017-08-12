package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
)

func main() {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	}))
	defer backendServer.Close()

	u, _ := url.Parse(backendServer.URL)
	http.Handle("/", httputil.NewSingleHostReverseProxy(u))

	log.Fatal(http.ListenAndServe(":9090", nil))
}
