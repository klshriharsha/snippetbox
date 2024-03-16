package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(staticFileSystem{http.Dir("./ui/static/")})
	// file server looks for the file under `./ui/static/`
	// so strip the `/static` prefix from request URL
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("listening on port 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
