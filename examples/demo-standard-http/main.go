package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/alexisvisco/kcd"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	r.Get("/{name}", kcd.Handler(StandardHttpHandler, http.StatusOK))
	_ = http.ListenAndServe(":3000", r)
}

// StandardHttpHandler is a standard http handler and work with kcd.
// The default status code does not works since you are controlling
// the return via http.ResponseWriter.
//
// This maybe useful in some cases:
// - if you have a large codebase and you want to progressively
//   integrate kcd
// - if you have a complex code for instance websocket, sse ...
func StandardHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
