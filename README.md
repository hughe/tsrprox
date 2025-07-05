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

TODO: document the command line interface.

