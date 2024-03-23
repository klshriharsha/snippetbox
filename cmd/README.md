# cmd

This directory contains only application-specific executables.

### web

This is the web server that can be executed by running `go run ./cmd/web` from the root directory

### Running web server

-   Generate a self-signed TLS cert in order to run the server in HTTPS mode. Use the command `go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost`. Place the cert and key in `/tls`
-   Use `go run ./cmd/web >>./info.log 2>>./error.log` to redirect logs to files on disk
-   Use `go run ./cmd/web -help` for documentation on accepted commandline flags
