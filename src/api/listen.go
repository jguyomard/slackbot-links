package api

import (
	"fmt"
	"net/http"

	"../config"
	"../links"
)

func Listen() {
	conf := config.Get()

	// Routes
	http.HandleFunc("/links/", handleLinks)

	// Listen!
	fmt.Println("API Listen on port", conf.APIListenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.APIListenPort), nil)
}

func handleLinks(w http.ResponseWriter, r *http.Request) {
	links := links.Search(r.URL.Query())
	fmt.Fprintf(w, collection(links, linkTransformer))
}
