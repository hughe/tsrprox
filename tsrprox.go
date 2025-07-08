package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
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
var authKeyFile string
var stateDir string

func init() {
	flag.IntVar(&httpPort, "http-port", 80, "The port on name to listen on for HTTP connections.  See -name. ")
	flag.BoolVar(&useHttp, "http", true, "Proxy HTTP connections on the port specified by -http-port.  -http=false to turn off.")
	flag.IntVar(&httpsPort, "https-port", 443, "The port on name to listen on for HTTP connections.  See -name. ")
	flag.BoolVar(&useHttps, "https", false, "Proxy HTTPS connections on the port specified by -https-port. "+
		"A TLS certificate will be provisioned from the tailnet.  "+
		"Note: to avoid certificate warnings you will need to used the fully qualified "+
		"domain name of the service.  E.g, https://name.tailXXXX.ts.net")

	flag.StringVar(&serverName, "name", "", "The DNS name of server on the tailnet, (defaults to the binary name)") // TODO: should this be an arg and should it default?
	flag.StringVar(&target, "target", "http://localhost", "The target URL to proxy requests to.")                   // TODO: should this be an Arg?
	flag.StringVar(&authKeyFile, "auth-key-file", "", "Path to file containing the Tailscale auth key")
	flag.StringVar(&stateDir, "state-dir", "",
		"Directory to store tsnet state files "+
			"(defaults to tsnet's default, see https://tailscale.com/kb/1522/tsnet-server#serverdir)")
}

var localClient *tailscale.LocalClient

func main() {
	var wg sync.WaitGroup
	flag.Parse()

	s := &tsnet.Server{}

	if serverName != "" {
		s.Hostname = serverName
	}

	if stateDir != "" {
		s.Dir = stateDir
	}

	if authKeyFile != "" {
		authKeyBytes, err := ioutil.ReadFile(authKeyFile)
		if err != nil {
			log.Fatal("Failed to read auth key file:", err)
		}
		s.AuthKey = strings.TrimSpace(string(authKeyBytes))
	}

	t, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Start()
	if err != nil {
		log.Fatal("Error starting tsnet server:", err)
	}
	defer s.Close()

	serveProxy := func(ln net.Listener, proto string, port int) {
		rp := httputil.NewSingleHostReverseProxy(t)

		log.Printf("Serving %s on port %d", proto, port)
		err = http.Serve(ln, rp)

		if err != nil {
			log.Printf("Error in %s on port %d: %s", proto, port, err)
		}
		wg.Done()
	}

	var httpListener, httpsListener net.Listener

	if useHttp {
		addr := fmt.Sprintf(":%d", httpPort)
		httpListener, err = s.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		defer httpListener.Close()

		wg.Add(1)
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

		wg.Add(1)
		go serveProxy(httpsListener, "HTTPS", httpsPort)
	}

	wg.Wait()
}

// Function director takes a HTTP request and decorates it with a
// header identifying the identity of the request sender.
//
// Not used right now.
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
