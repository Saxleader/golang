package picture

import (
	"html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("assets/root.html"))
	t.Execute(rw, nil)
}
