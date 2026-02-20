package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
    ##         .
    ## ## ## ==
    ## ## ## ## ## ===
   /"""""""""""""""""\___/ ===
  {                       /  ===-
   \______ O           __/
     \    \         __/
      \____\_______/

Hello from Docker!
`)
}

func main() {
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "80"
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handler)
	log.Printf("Go backend listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
