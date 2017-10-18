package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/jguyomard/slackbot-links/src/config"
	"github.com/jguyomard/slackbot-links/src/links"
)

// Listen starts api web server
func Listen() {
	conf := config.Get()

	// Routes
	router := httprouter.New()
	router.GET("/v1/links", handleSearchLinks)
	router.GET("/v1/links/:id", handleGetLink)
	router.DELETE("/v1/links/:id", handleDeleteLink)
	router.NotFound = http.HandlerFunc(handleErrorNotFound)

	// Listen!
	fmt.Println("API Listen on port", conf.APIListenPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.APIListenPort), NewAPIMiddleware(router)); err != nil {
		panic(err)
	}
}

func handleSearchLinks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	linksResult, err := links.Search(r.URL.Query())
	if err != nil {
		http.Error(w, errors(err.Error(), 400).ToJSON(), 400)
		return
	}
	meta := map[string]interface{}{
		"cursor": linksResult.GetCursor(),
		"total":  linksResult.GetTotal(),
	}
	fmt.Fprintf(w, "%s", collection(linksResult.GetLinks(), linkTransformer).SetMeta(meta).ToJSON())
}

func handleGetLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	link, found := links.Get(ps.ByName("id"))
	if !found {
		http.Error(w, errors("Link not found", 404).ToJSON(), 404)
		return
	}

	fmt.Fprintf(w, "%s", item(link, linkTransformer).ToJSON())
}

func handleDeleteLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	link, found := links.Get(ps.ByName("id"))
	if !found {
		http.Error(w, errors("Link not found", 404).ToJSON(), 404)
		return
	}

	// Delete this link
	if !link.Delete() {
		http.Error(w, errors("Unable to delete this link", 500).ToJSON(), 500)
		return
	}

	w.WriteHeader(204)
}

func handleErrorNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, errors("Route not found", 404).ToJSON(), 404)
}
