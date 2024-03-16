# cmd

This directory contains only application-specific executables.

### web

This is the web server that can be executed by running `go run ./cmd/web` from the root directory

### Running web server

-   Use `go run ./cmd/web >>./info.log 2>>./error.log` to redirect logs to files on disk
-   Use `go run ./cmd/web -help` for documentation on accepted commandline flags
