package googleapi

import (
	// "encoding/json"
	// "fmt"
	"html/template"
	// "io/ioutil"
	"log"
	"net/http"
	// "net/url"
	"google.golang.org/api/gmail/v1"
)

type resultinput struct {
	Search string
	Emails []string
}

func init() {
	http.Handle("/hiddenDir/", http.StripPrefix("/hiddenDir/", http.FileServer(http.Dir("assets/"))))
	http.HandleFunc("/", apiQuery)
	http.HandleFunc("/result", results)
}

func apiQuery(rw http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("assets/apiQuery.html"))
	err := t.ExecuteTemplate(rw, "apiQuery", nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func results(rw http.ResponseWriter, req *http.Request) {
	var mysearch = resultinput{Search: req.FormValue("search")}

	// safeAddr := url.QueryEscape(mysearch.Search)
	// fullURL := fmt.Sprintf("https://www.googleapis.com/gmail/v1/users/me/messages?q=%s", mysearch)

	client := &http.Client{}

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
	if err != nil {
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

	for _, value := range r.Messages {
		mysearch.Emails = append(mysearch.Emails, "https://mail.google.com/mail/u/0/#all/"+value.Id)
	}

	t := template.Must(template.ParseFiles("assets/results.html"))
	err = t.ExecuteTemplate(rw, "results", mysearch)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
