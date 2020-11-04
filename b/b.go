package main

import (
	"fmt"
	"net/http"
	"log"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

type Result struct {
	Status string
}

func main() {

	http.HandleFunc("/", home)
	http.ListenAndServe(":9091", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("ccNumber")

	resultCoupon := makeHttpCall("http://localhost:9092", coupon)

	result := Result{Status: "declined"}

	if ccNumber == "1" {
		result.Status = "approved"
	}

	if resultCoupon.Status == "invalid" {
		result.Status = "invalid coupon"
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error Processing Json")
	}

	fmt.Fprintf(w, string(jsonData))

}

func makeHttpCall(urlMicroservice string, coupon string) Result {
	values := url.Values{}
	values.Add("coupon", coupon)

	res, err := http.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Status: "Servidor c fora do ar"}
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