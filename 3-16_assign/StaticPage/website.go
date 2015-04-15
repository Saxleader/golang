package staticpage

import (
	"html/template"
	"log"
	"net/http"
)

func init() {
	http.Handle("/hiddenDir/", http.StripPrefix("/hiddenDir/", http.FileServer(http.Dir("assets/"))))
	http.HandleFunc("/", handler)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("index.html"))
	err := t.ExecuteTemplate(rw, "index.html", nil)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
}
