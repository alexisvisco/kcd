package main

import (
	"net/http"

	"github.com/alexisvisco/kcd"
)

func main() {

	//                                                                    v   status code will be ignored
	http.HandleFunc("/example/", kcd.Handler(StandardHttpHandler, http.StatusOK))
	http.ListenAndServe(":8080", nil)
}

// StandardHttpHandler is a standard http handler and work with kcd.
// The default status code does not works since you are controlling
// the return via http.ResponseWriter.
//
// This maybe useful in some cases:
// - if you have a large codebase and you want to progressively
//   integrate kcd
// - if you have a complex code for instance websocket, sse ...
//
//                                                              v as you can see there is no error, and no output struct
func StandardHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
