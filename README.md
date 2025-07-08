# tsrprox

A very simple [Tailscale](https://tailscale.com/) reverse proxy. 

Works like the [`tailscale
serve`](https://tailscale.com/kb/1242/tailscale-serve) command except
that you can give each service a different on the network.  For
example if I have a host `server.tailXXXX.ts.net` and you want to
expose a service running on that host with a different memorable name
e.g., `myservice.tailXXXX.ts.net`.


Inspired by [this entry in the Tailscale
blog](https://tailscale.com/blog/tsnet-virtual-private-services).

## Build

**TODO**

## Run

```
./tsrprox -help
Usage of ./tsrprox:
  -auth-key-file string
    	Path to file containing the Tailscale auth key
  -http
    	Proxy HTTP connections on the port specified by -http-port.  -http=false to turn off. (default true)
  -http-port int
    	The port on name to listen on for HTTP connections.  See -name.  (default 80)
  -https
    	Proxy HTTPS connections on the port specified by -https-port. A TLS certificate will be provisioned from the tailnet.  Note: to avoid certificate warnings you will need to used the fully qualified domain name of the service.  E.g, https://name.tailXXXX.ts.net
  -https-port int
    	The port on name to listen on for HTTP connections.  See -name.  (default 443)
  -name string
    	The DNS name of server on the tailnet, (defaults to the binary name)
  -state-dir string
    	Directory to store tsnet state files (defaults to tsnet's default, see https://tailscale.com/kb/1522/tsnet-server#serverdir)
  -target string
    	The target URL to proxy requests to. (default "http://localhost")
```

### Examples 

Expose the service running on `localhost` port `8080` as `bar` on port
`80` in your tailnet.  After this the service will be accessible as
`http://bar` or `http://bar.tailXXXX.ts.net` within your tailnet.

```
$ export TS_AUTHKEY=tskey-auth-XXXXXXX-XXXXXXXXXXXXXXXXXXXXX
$ tsrprox -target http://localhost:8080 -name bar
```

Expose the service running on `localhost` port `8000` as `baz` or
`baz.tailXXXX.ts.net` on via HTTPS and HTTP.  After this the service
will be accessible as `https://baz.tailXXXX.ts.net`, `http://baz` or
`http://baz.tailXXXX.ts.net` within your tailnet.  

```
$ export TS_AUTHKEY=tskey-auth-XXXXXXX-XXXXXXXXXXXXXXXXXXXXX
$ tsrprox -target http://localhost:8000 -name bar -https
```

Note that to avoid HTTPS errors like "This Connection Is Not Private"
you will need to use the fully qualified domain name (FQDN) with
HTTPS.  E.g., `https://baz.tailXXXX.ts.net` will work, `https://baz`
will not work.  This is because the services TLS certificate only
contains the FQDN.




