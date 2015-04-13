package urldata

import (
	"html/template"
	"log"
	"net/http"
)

var mytmpl *template.Template

func init() {
	mytmpl = template.Must(template.ParseFiles("assets/root.html", "assets/confirm.html", "assets/welcome.html", "assets/redirect.html"))
	http.HandleFunc("/", myroot)
	http.HandleFunc("/correct", confirm)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/welcome", greeting)
}

func greeting(rw http.ResponseWriter, req *http.Request) {
	data := struct {
		Name string
		Age  string
	}{req.FormValue("name"), req.FormValue("age")}
	err := mytmpl.ExecuteTemplate(rw, "welcome.html", data)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
}

func queryform(rw http.ResponseWriter, data interface{}) {
	tmpl := "root.html"
	// fmt.Fprintf(rw, "Data: %v", data)
	if data != nil {
		tmpl = "redirect.html"
	}
	err := mytmpl.ExecuteTemplate(rw, tmpl, data)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
}

func myroot(rw http.ResponseWriter, req *http.Request) {
	// mytmpl = template.Must(template.ParseFiles("assets/root.html", "assets/confirm.html", "assets/welcome.html", "assets/redirect.html"))
	// data := struct {
	// 	Name string
	// 	Age  string
	// }{"name", "age"}
	queryform(rw, nil)
	// err := mytmpl.ExecuteTemplate(rw, "root.html", nil)
	// if err != nil {
	// 	log.Fatalf("Error: %v", err.Error())
	// }
}

func confirm(rw http.ResponseWriter, req *http.Request) {
	data := struct {
		Name string
		Age  string
	}{req.FormValue("name"), req.FormValue("age")}
	err := mytmpl.ExecuteTemplate(rw, "confirm.html", data)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
}

func edit(rw http.ResponseWriter, req *http.Request) {
	data := struct {
		Name string
		Age  string
	}{req.FormValue("name"), req.FormValue("age")}
	queryform(rw, data)
}
