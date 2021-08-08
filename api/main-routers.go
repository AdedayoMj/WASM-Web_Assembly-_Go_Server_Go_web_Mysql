package api

import (
	"fmt"
	"gowebapp/config"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

// MainRouters are the collection of all URLs for the Main App.
func MainRouters(r *mux.Router) {
	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc("/articles/{category}/", ArticlesCategoryHandler)
	r.HandleFunc("/products/{product_name}/{id:[0-9]+}/", ProductInfo)
}

// contextData are the most widely use common variables for each pages to load.
type contextData map[string]interface{}

// Home function is to render the homepage page.
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(config.SiteRootTemplate+"front/index.html", config.SiteHeaderTemplate, config.SiteFooterTemplate))

	data := contextData{
		"PageTitle":    "Welcome to Adedayo GoLang Test",
		"PageMetaDesc": config.SiteSlogan,
		"CanonicalURL": r.RequestURI,
		"CsrfToken":    csrf.Token(r),
		"Settings":     config.SiteSettings,
	}
	tmpl.Execute(w, data)
}

//Article Category Handler ...
func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

//ProductInfo ...
func ProductInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Product Name: %v\n", vars["product_name"])
	fmt.Fprintf(w, "Product ID: %v\n", vars["id"])
}
