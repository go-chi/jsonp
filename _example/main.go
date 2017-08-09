// JSONP example using Chi http router.. but anything that accepts
// a http.Handler will work
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jsonp"
	"github.com/go-chi/render"
)

func main() {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(jsonp.Handler)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := &SomeObj{"superman"}
		render.JSON(w, r, data)
	})

	err := http.ListenAndServe(":4444", mux)
	if err != nil {
		log.Fatal(err)
	}
}

type SomeObj struct {
	Name string `json:"name"`
}
