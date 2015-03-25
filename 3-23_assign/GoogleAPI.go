package googleAPI

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type resultinput struct {
	Search string
	Emails []string
}

func init() {
	http.HandleFunc("/", apiQuery)
	http.HandleFunc("/result", results)
}

func apiQuery(rw http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("assets/apiQuery"))
	err := t.Execute(rw, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func results(rw http.ResponseWriter, req *http.Request) {
	var mysearch = resultinput{Search: req.FormValue("search")}

	safeAddr := url.QueryEscape(mysearch.Search)
	fullURL := fmt.Sprintf("https://www.googleapis.com/gmail/v1/users/me/messages?q=%s", mysearch)

	client := &http.Client{}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	resp, requestErr := client.Do(req)
	if requestErr != nil {
		log.Fatal("Do: ", requestErr)
	}

	defer resp.Body.Close()

	body, dataReadErr := ioutil.ReadAll(resp.Body)
	if dataReadErr != nil {
		log.Fatal("ReadAll: ", dataReadErr)
	}

	res := make(map[string][]map[string]interface{}, 0)

	json.Unmarshal(body, &res)

	var temp string

	for i, _ := range res["messages"] {
		temp, _ = res["messages"][i]["id"]
		append(mysearch.Emails, temp)
	}

	var myemails []string

	for i, value := range mysearch.Emails {
		mysearch.Emails[i] = "https://mail.google.com/mail/u/0/#inbox/" + value
	}

	t := template.Must(template.ParseFiles("assets/results.html"))
	err = t.Execute(rw, mysearch)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
