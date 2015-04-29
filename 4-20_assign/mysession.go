package mysessions

import (
	"html/template"
	"net/http"

	"github.com/astaxie/beego/session"
)

var globalSessions *session.Manager

var mytmpl *template.Template

type userdata struct {
	User    string
	Message string
}

func init() {
	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 3600, "providerConfig": ""}`)
	go globalSessions.GC()
	http.Handle("/hiddenDir/", http.StripPrefix("/hiddenDir/", http.FileServer(http.Dir("assets/"))))
	http.HandleFunc("/", myroot)
	http.HandleFunc("/store", store)
	http.HandleFunc("/erase", erase)
	mytmpl = template.Must(template.ParseFiles("assets/myroot.html"))
}

func myroot(rw http.ResponseWriter, req *http.Request) {
	sess, _ := globalSessions.SessionStart(rw, req)
	var mydata userdata
	tempdata := sess.Get("userdata")
	if tempdata != nil {
		mydata = tempdata.(userdata)
	}
	mytmpl.ExecuteTemplate(rw, "myroot.html", mydata)
}

func store(rw http.ResponseWriter, req *http.Request) {
	sess, _ := globalSessions.SessionStart(rw, req)
	var mydata userdata
	mydata.User = req.FormValue("name")
	mydata.Message = req.FormValue("memo")
	sess.Set("userdata", mydata)
	http.Redirect(rw, req, "/", http.StatusFound)
}

func erase(rw http.ResponseWriter, req *http.Request) {
	sess, _ := globalSessions.SessionStart(rw, req)
	userdata := sess.Get("userdata")
	if userdata != nil {
		globalSessions.SessionDestroy(rw, req)
	}
	http.Redirect(rw, req, "/", http.StatusFound)
}
