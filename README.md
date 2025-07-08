# tsrprox

A simple [Tailscale](https://tailscale.com/) reverse proxy that allows
you to expose local services with custom DNS names on your Tailnet.

Works like the [`tailscale
serve`](https://tailscale.com/kb/1242/tailscale-serve) command except
that you can give each service a different name on the network. For
example, if you have a host `server.tailXXXX.ts.net` and want to
expose a service running on that host with a different memorable name
like `myservice.tailXXXX.ts.net`.


Inspired by [this entry in the Tailscale
blog](https://tailscale.com/blog/tsnet-virtual-private-services).

## Prerequisites

- Go 1.24 or later
- A Tailscale account and auth key
    + The auth key should be reusable and ephemeral.

## Build

First, install dependencies:
```bash
go mod download
```

Then build:
```bash
# Build for current platform
make

# Build for Linux (amd64)
make tsrprox.linux

# Direct Go build
go build .
```

## Run

```bash
./tsrprox -help
Usage of ./tsrprox:
  -auth-key-file string
    	Path to file containing the Tailscale auth key
  -http
    	Proxy HTTP connections on the port specified by -http-port.  -http=false to turn off. (default true)
  -http-port int
    	The port on name to listen on for HTTP connections.  See -name.  (default 80)
  -https
    	Proxy HTTPS connections on the port specified by -https-port. A TLS certificate will be provisioned 
		from the tailnet.  Note: to avoid certificate warnings you will need to used the fully qualified
		domain name of the service.  E.g, https://name.tailXXXX.ts.net
  -https-port int
    	The port on name to listen on for HTTP connections.  See -name.  (default 443)
  -name string
    	The DNS name of server on the tailnet, (defaults to the binary name)
  -state-dir string
    	Directory to store tsnet state files (defaults to tsnet's default, see 
		https://tailscale.com/kb/1522/tsnet-server#serverdir)
  -target string
    	The target URL to proxy requests to. (default "http://localhost")
```

### Examples 

**Basic HTTP proxy:**

Expose the service running on `localhost` port `8080` as `bar` on port
`80` in your tailnet. After this the service will be accessible as
`http://bar` or `http://bar.tailXXXX.ts.net` within your tailnet.

```bash
$ export TS_AUTHKEY=tskey-auth-XXXXXXX-XXXXXXXXXXXXXXXXXXXXX
$ tsrprox -target http://localhost:8080 -name bar
```

**HTTPS and HTTP proxy:**

Expose the service running on `localhost` port `8000` as `baz` via
HTTPS and HTTP. After this the service will be accessible as
`https://baz.tailXXXX.ts.net`, `http://baz` or
`http://baz.tailXXXX.ts.net` within your tailnet.

```bash
$ export TS_AUTHKEY=tskey-auth-XXXXXXX-XXXXXXXXXXXXXXXXXXXXX
$ tsrprox -target http://localhost:8000 -name baz -https
```

**Using auth key file:**

```bash
$ echo "tskey-auth-XXXXXXX-XXXXXXXXXXXXXXXXXXXXX" > /tmp/authkey
$ tsrprox -target http://localhost:3000 -name myapp -auth-key-file /tmp/authkey
```

## Authentication

You need a Tailscale auth key to use this tool. You can:
1. Set the `TS_AUTHKEY` environment variable
2. Use the `-auth-key-file` flag to point to a file containing your auth key

Get your auth key from the [Tailscale admin
console](https://login.tailscale.com/admin/settings/keys).

## Troubleshooting

**HTTPS Certificate Warnings:**

To avoid HTTPS errors like "This Connection Is Not Private", you must
use the fully qualified domain name (FQDN) with HTTPS. For example:
- ✅ `https://baz.tailXXXX.ts.net` will work
- ❌ `https://baz` will not work

This is because the TLS certificate only contains the FQDN of the proxy.

**Connection Issues:**

- Ensure your Tailscale client is running and connected
- Verify your auth key is valid and has the necessary permissions
- Check that the target service is accessible locally
