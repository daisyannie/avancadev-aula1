package main

import (
	"html/template"
	"net/http"
	"log"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
)

type Result struct {
	Status string
}

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.ListenAndServe(":9090", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/home.html"))
	t.Execute(w, Result{})
}

func process(w http.ResponseWriter, r *http.Request) {

	result := makeHttpCall("http://localhost:9091", r.FormValue("coupon"), r.FormValue("cc-number"))
	t := template.Must(template.ParseFiles("templates/home.html"))
	t.Execute(w, result)
}

func makeHttpCall(urlMicroservice string, coupon string, ccNumber string) Result {
	values := url.Values{}
	values.Add("coupon", coupon)
	values.Add("ccNumber", ccNumber)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	
	res, err := retryClient.PostForm(urlMicroservice, values)
	if err != nil {
		//log.Fatal("Microservice b out")
		result := Result{Status: "Servidor b fora do ar"}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Erro ao processar resultado")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result
}