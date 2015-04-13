package cookies

import (
	"html/template"
	"log"
	"net/http"
)

var mytmpl *template.Template

func init() {
	mytmpl = template.Must(template.ParseFiles("assets/root.html", "assets/seeCookie.html", "assets/setCookie.html"))
	http.HandleFunc("/", handler)
	http.HandleFunc("/setCookie", setMyCookie)
	http.HandleFunc("/seeCookie", seeMyCookie)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	err := mytmpl.ExecuteTemplate(rw, "root.html", nil)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
}

func setMyCookie(rw http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	age := req.FormValue("age")
	namecookie := http.Cookie{Name: "name", Value: name, Path: "/seeCookie"}
	agecookie := http.Cookie{Name: "age", Value: age, Path: "/seeCookie"}
	http.SetCookie(rw, &namecookie)
	http.SetCookie(rw, &agecookie)
	err := mytmpl.ExecuteTemplate(rw, "setCookie.html", nil)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
	// fmt.Fprintf(rw, "name: %v age: %v", name, age)
	// http.Redirect(rw, req, "/seeCookie", http.StatusFound)
	// http.Redirect(w, r, urlStr, code)
}

func seeMyCookie(rw http.ResponseWriter, req *http.Request) {
	nameCookie, err := req.Cookie("name")
	var name, age string
	if err != nil {
		log.Printf("Error: %v", err.Error())
	} else {
		name = nameCookie.Value
	}
	ageCookie, err := req.Cookie("age")
	if err != nil {
		log.Printf("Error: %v", err.Error())
	} else {
		age = ageCookie.Value
	}
	err = mytmpl.ExecuteTemplate(rw, "seeCookie.html", struct {
		Name string
		Age  string
	}{name, age})
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
}
