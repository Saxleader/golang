package googleapi

import (
	// "encoding/json"
	// "fmt"
	"encoding/json"
	"html/template"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	// "google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
)

// Redirect URIs https://curious-cistern-90523.appspot.com/results
// Javascript Origins https://curious-cistern-90523.appspot.com
var conf = &oauth2.Config{
	ClientID:     "506452277892-0cqothcekkv904em09410b52iqngbd74.apps.googleusercontent.com",
	ClientSecret: "hRsmnCOhUnWvWieKNltXKXkT",
	RedirectURL:  "https://curious-cistern-90523.appspot.com/results",
	Scopes:       []string{gmail.GmailReadonlyScope},
	Endpoint:     google.Endpoint,
}

var mytmpl *template.Template

type resultinput struct {
	ID      string
	Snippet string
}

func init() {
	mytmpl = template.Must(template.ParseFiles("assets/apiQuery.html", "assets/results.html"))
	http.Handle("/hiddenDir/", http.StripPrefix("/hiddenDir/", http.FileServer(http.Dir("assets/"))))
	http.HandleFunc("/", apiQuery)
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/results", results)
}

func apiQuery(rw http.ResponseWriter, req *http.Request) {
	err := mytmpl.ExecuteTemplate(rw, "apiQuery", nil)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func submit(rw http.ResponseWriter, req *http.Request) {
	url := conf.AuthCodeURL(req.FormValue("search"), oauth2.AccessTypeOnline)
	http.Redirect(rw, req, url, http.StatusFound)
}

func results(rw http.ResponseWriter, req *http.Request) {
	search := req.FormValue("state")
	code := req.FormValue("code")

	c := appengine.NewContext(req)

	mytoken, err := conf.Exchange(c, code)
	// if err != nil {
	// 	log.Printf("Error: %v", err.Error())
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// }

	client := urlfetch.Client(c)

	resp, err := client.Get("https://www.googleapis.com/gmail/v1/users/me/messages?q=" + search + "&access_token=" + mytoken.AccessToken)
	if err != nil {
		log.Printf("Client: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Printf("ReadAll: %v", err.Error())
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// }

	var mymail gmail.ListMessagesResponse

	err = json.Unmarshal(body, &mymail)
	if err != nil {
		log.Printf("JSON: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var IDlist []resultinput

	for _, val := range mymail.Messages {
		IDlist = append(IDlist, resultinput{val.Id, val.Snippet})
	}

	// safeAddr := url.QueryEscape(mysearch.Search)
	// fullURL := fmt.Sprintf("https://www.googleapis.com/gmail/v1/users/me/messages?q=%s", mysearch)
	// u, err := user.CurrentOAuth(c, "https://www.googleapis.com/auth/gmail.readonly")
	// if err != nil {
	// 	log.Fatalf("Unable to authenticate user: %v", err)
	// }

	// client := urlfetch.Client(c)

	// svc, err := gmail.New(client)
	// if err != nil {
	// 	log.Fatalf("Unable to create gmail service: %v", err)
	// }

	// res := svc.Users.Messages.List("me").Q(mysearch.Search)

	// r, err := res.Do()

	// req, err := http.NewRequest("GET", fullURL, nil)
	// if err != nil {
	// 	log.Fatal("NewRequest: ", err)
	// }

	// resp, requestErr := client.Do(req)
	/*if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	// defer resp.Body.Close()

	// body, dataReadErr := ioutil.ReadAll(resp.Body)
	// if dataReadErr != nil {
	// 	log.Fatal("ReadAll: ", dataReadErr)
	// }

	// res := make(map[string][]map[string]interface{}, 0)

	// json.Unmarshal(body, &res)

	// for i, _ := range res["messages"] {
	// 	temp, _ = res["messages"][i]["id"]
	// 	append(mysearch.Emails, temp)
	// }

	// var myemails []string

	/*for _, value := range r.Messages {
		mysearch.Emails = append(mysearch.Emails, "https://mail.google.com/mail/u/0/#all/"+value.Id)
	}*/

	err = mytmpl.ExecuteTemplate(rw, "results", IDlist)
	if err != nil {
		// fmt.Fprintf(rw, "An error has occured again: SEARCH:%v CLIENT:%v SVC:%v RES:%v R:%v ERR1:%v", mysearch.Search, client, svc, res, r, err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
