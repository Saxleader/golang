package googleapi

import (
	// "encoding/json"
	"fmt"
	"html/template"
	"log"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
)

// Redirect URIs https://curious-cistern-90523.appspot.com/oauth2callback
// Javascript Origins https://curious-cistern-90523.appspot.com
var conf = &oauth2.Config{
	ClientID:	"506452277892-0cqothcekkv904em09410b52iqngbd74.apps.googleusercontent.com"
	ClientSecret:	"hRsmnCOhUnWvWieKNltXKXkT"
	Scopes:	gmail.GmailReadonlyScope
	Endpoint: google.Endpoint,
}

var mytmpl *template.Template

type resultinput struct {
	Search string
	Emails []string
}

func init() {
	mytmpl = template.Must(template.ParseFiles("assets/apiQuery.html","assets/results.html"))
	http.Handle("/hiddenDir/", http.StripPrefix("/hiddenDir/", http.FileServer(http.Dir("assets/"))))
	http.HandleFunc("/", apiQuery)
	http.HandleFunc("/result", results)
}

func apiQuery(rw http.ResponseWriter, req *http.Request) {
	err := mytmpl.ExecuteTemplate(rw, "apiQuery", nil)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func results(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	search := req.FormValue("search")
	client := conf.Client(c, t)
	resp, err := client.Get("https://www.googleapis.com/gmail/v1/me/userId/messages?q="+search)
	if err != nil {
		log.Printf("Client: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var mymail gmail.ListMessagesResponse

	err = json.Unmarshal(body, &mymail)
	if err != nil {
		log.Printf("JSON: %v", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	// safeAddr := url.QueryEscape(mysearch.Search)
	// fullURL := fmt.Sprintf("https://www.googleapis.com/gmail/v1/users/me/messages?q=%s", mysearch)
	// u, err := user.CurrentOAuth(c, "https://www.googleapis.com/auth/gmail.readonly")
	// if err != nil {
	// 	log.Fatalf("Unable to authenticate user: %v", err)
	// }

	client := urlfetch.Client(c)

	svc, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to create gmail service: %v", err)
	}

	res := svc.Users.Messages.List("me").Q(mysearch.Search)

	r, err := res.Do()

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

	// t := template.Must(template.ParseFiles("assets/results.html"))
	// err = t.ExecuteTemplate(rw, "results", mysearch)
	// if err != nil {
	fmt.Fprintf(rw, "An error has occured again: SEARCH:%v CLIENT:%v SVC:%v RES:%v R:%v ERR1:%v", mysearch.Search, client, svc, res, r, err)
	// http.Error(rw, err.Error(), http.StatusInternalServerError)

	// }
}
