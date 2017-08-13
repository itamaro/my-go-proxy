package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Prox struct {
	target     *url.URL
	proxy      *httputil.ReverseProxy
	numProxies int
}

func New(target string, numProxies int) *Prox {
	url, _ := url.Parse(target)

	return &Prox{
		target:     url,
		proxy:      httputil.NewSingleHostReverseProxy(url),
		numProxies: numProxies,
	}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Remote Address: %s\n", GetRemoteAddr(r, p.numProxies))
	fmt.Printf("Request URI: %s\n", r.RequestURI)
	requestDump, _ := httputil.DumpRequest(r, true)
	// if err != nil {
	//   fmt.Println(err)
	// }
	fmt.Println(string(requestDump))
	p.proxy.ServeHTTP(w, r)
}

func main() {
	const (
		defaultPort            = ":80"
		defaultPortUsage       = "proxy port to listen on"
		defaultTarget          = "http://127.0.0.1:8080"
		defaultTargetUsage     = "target URL to redirect to"
		defaultNumProxies      = -1
		defaultNumProxiesUsage = "how many proxies are expected before this one"
	)

	// flags
	port := flag.String("port", defaultPort, defaultPortUsage)
	url := flag.String("url", defaultTarget, defaultTargetUsage)
	numProxies := flag.Int("num-proxies", defaultNumProxies, defaultNumProxiesUsage)

	flag.Parse()

	fmt.Printf("proxy listening on: %s\n", *port)
	fmt.Printf("reverse-proxying to: %s\n", *url)

	// proxy
	proxy := New(*url, *numProxies)

	// server
	http.HandleFunc("/", proxy.handle)
	log.Fatal(http.ListenAndServe(*port, nil))
}
