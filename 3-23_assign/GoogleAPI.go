package googleAPI

import (
	"html/template"
	"net/http"
	"net/url"
	"fmt"
)

func init() {
	http.HandleFunc("/", apiQuery)
	http.HandleFunc("/result", results)
}

func apiQuery(rw http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("assets/apiQuery"))
	err = t.Execute(rw, nil)
	if err != nil {
		http.Error(rw, err.error(), http.StatusInternalServerError)
	}
}

func results(rw http.ResponseWriter, req *http.Request) {
	mysearch := req.FormValue("search")

	safeAddr := url.QueryEscape(mysearch)
	fullURL := fmt.Sprintf("https://www.googleapis.com/gmail/v1/users/me/messages?q=", ...)

	t := template.Must(template.ParseFiles("assets/results.html"))
	err = t.Execute(rw, data)
	if err != nil {
		http.Error(rw, err.error(), http.StatusInternalServerError)
	}
}
