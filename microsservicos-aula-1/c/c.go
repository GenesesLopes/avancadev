package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/hashicorp/go-retryablehttp"
)

type Coupons struct {
	Coupon []string
}

func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item {
			return "valid"
		}
	}
	return "invalid"
}

type Result struct {
	Status string
}

var coupons Coupons

func main() {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	response, err := retryClient.Get("http://php_dev:8000")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer response.Body.Close()
	
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(responseData,&coupons.Coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	valid := coupons.Check(coupon)

	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))

}