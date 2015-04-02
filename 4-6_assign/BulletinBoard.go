package bulletinboard

import (
	"net/http"
	"text/template"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

type Post struct {
	Author     string
	Message    string
	UpdateDate string
	PostDate   string
	Updated    bool
}

func init() {
	http.Handle("/hiddenDir/", http.StripPrefix("/hiddenDir/", http.FileServer(http.Dir("assets/"))))
	http.HandleFunc("/", board)
	http.HandleFunc("/view/", viewpost)
	http.HandleFunc("/edit", editpost)
}

// The following handlerfunc was my first option for login and authentication.
// I decided to go with the authentication in the app.yaml.
/*func board(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-type", "text/html")
	c := appengine.NewContext(req)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, req.URL.String())
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Location", url)
		rw.WriteHeader(http.StatusFound)
		return
	}
	t := template.Must(template.ParseFiles("assets/board.html"))
	t.ExecuteTemplate(wr, "Board", nil)
}*/

func board(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-type", "text/html")
	c := appengine.NewContext(req)
	u := user.Current(c)
	rw.Header().Set("Location", req.URL.String())
	rw.WriteHeader(http.StatusFound)
	t := template.Must(template.ParseFiles("assets/board.html"))
	t.ExecuteTemplate(rw, "Board", u)
}

func viewpost(rw http.ResponseWriter, req *http.Request) {

}

func editpost(rw http.ResponseWriter, req *http.Request) {

}

func updatepost(rw http.ResponseWriter, req *http.Request) {

}

func createpost(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	u := user.Current(c)
	m := req.FormValue("Message")
	t := time.Now().Format("Jan 2, 2006 3:04 PM")
	var p = Post{Author: u.String(), Message: m, PostDate: t, Updated: false}
	datastore.Put(c, datastore.NewIncompleteKey(c, "Post", nil), &p)
}
