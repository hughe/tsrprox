# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a Go-based Tailscale reverse proxy server that allows exposing local services through Tailscale with custom DNS names. It's a simple alternative to `tailscale serve` that provides more control over service naming on the Tailnet.

## Build Commands

- `make` or `make tsrprox` - Build the binary for current platform
- `make tsrprox.linux` - Cross-compile for Linux (amd64)
- `go build .` - Direct Go build command

## Architecture

The application is a single Go file (`tsrprox.go`) that:

1. **Main Components:**
   - Uses `tailscale.com/tsnet` for Tailscale network integration
   - Implements HTTP/HTTPS reverse proxy using `net/http/httputil`
   - Supports TLS certificate provisioning from Tailnet

2. **Key Configuration:**
   - `-name`: Custom DNS name for the service (defaults to binary name)
   - `-target`: Target URL to proxy requests to (default: http://localhost)
   - `-http`/`-https`: Enable HTTP (port 80) or HTTPS (port 443) listeners
   - `-http-port`/`-https-port`: Custom port configuration

3. **Network Flow:**
   - Creates tsnet.Server instance for Tailscale connectivity
   - Sets up HTTP/HTTPS listeners on specified ports
   - Proxies requests to target URL using single-host reverse proxy

## Development Notes

- Single-file Go application with minimal dependencies
- Uses Tailscale's tsnet library for network abstraction
- Incomplete features marked with TODO comments (identity forwarding, custom headers)
- No test files present in the codebase
- Uses Go 1.19 as specified in go.mod

## Runtime Behavior

The proxy serves requests concurrently using goroutines for each protocol (HTTP/HTTPS). The HTTPS listener automatically provisions TLS certificates through Tailscale's certificate management system.