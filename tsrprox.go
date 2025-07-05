package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"net/http"
	"net/http/httputil"
	"net/url"

	"tailscale.com/client/tailscale"
	"tailscale.com/tsnet"
)

var httpPort int
var useHttp bool
var httpsPort int
var useHttps bool
var serverName string
var target string
var forwardIdentity bool

func init() {
	flag.IntVar(&httpPort, "http-port", 80, "The port on name to listen on for HTTP connections.  See -name. ")
	flag.BoolVar(&useHttp, "http", true, "Proxy HTTP connections on the port specified by -http-port.  -http=false to turn off.")
	flag.IntVar(&httpPort, "https-port", 443, "The port on name to listen on for HTTP connections.  See -name. ")
	flag.BoolVar(&useHttp, "https", false, "Proxy HTTPS connections on the port specified by -https-port. "+
		"A TLS certificate will be provisioned from the tailnet.  "+
		"Note: to avoid certificate warnings you will need to used the fully qualified "+
		"domain name of the service.  E.g, https://name.tailXXXX.ts.net")

	flag.StringVar(&serverName, "name", "", "The DNS name of server on the tailnet, (defaults to the binary name)") // TODO: should this be an arg and should it default?
	flag.StringVar(&target, "target", "http://localhost", "The target URL to proxy requests to.")                   // TODO: should this be an Arg?
}

var localClient *tailscale.LocalClient

func main() {
	var wg sync.WaitGroup
	flag.Parse()

	s := &tsnet.Server{}

	if serverName != "" {
		s.Hostname = serverName
	}

	t, err := url.Parse(target)

	if err != nil {
		log.Fatal(err)
	}

	s.Start()
	defer s.Close()

	serveProxy := func(ln net.Listener, proto string, port int) {
		wg.Add(1)
		rp := httputil.NewSingleHostReverseProxy(t)

		log.Println("Serving %s on port %d", proto, port)
		err = http.Serve(ln, rp)

		if err != nil {
			log.Fatal(err)
		}
	}

	var httpListener, httpsListener net.Listener

	if useHttp {
		addr := fmt.Sprintf(":%d", httpPort)
		httpListener, err = s.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		defer httpListener.Close()

		go serveProxy(httpListener, "HTTP", httpPort)

	}

	if useHttps {
		localClient, err = s.LocalClient()
		if err != nil {
			log.Fatal(err)
		}

		addr := fmt.Sprintf(":%d", httpsPort)
		ln, err := s.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()

		httpsListener = tls.NewListener(ln, &tls.Config{
			GetCertificate: localClient.GetCertificate,
		})
		go serveProxy(httpsListener, "HTTPS", httpsPort)
	}

	wg.Wait()
}

func director(req *http.Request) {
	var id string = ""

	if forwardIdentity {
		if localClient != nil {
			if req.RemoteAddr != "" {
				res, err := localClient.WhoIs(context.TODO(), req.RemoteAddr)
				if err != nil {
					log.Fatal("Error getting identity")
				} else {
					id = res.UserProfile.LoginName
				}
			}
		}

		if id != "" {
			// TODO: set the identity header
		} else {
			// TODO: set an identity unavailable header
		}
	}

	// TODO: set proxied by header.

}
