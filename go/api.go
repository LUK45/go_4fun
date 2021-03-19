package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"regexp"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/gocolly/colly/v2"
	"github.com/gorilla/mux"
	"github.com/PuerkitoBio/goquery"
)

type Temperature struct {
	Date string `json:"date"`
	Low string `json:"low"`
	High string `json:"high"`
}

type TemperatureCity struct {
	Id string `json:"id"`
	City string `json:"city"`
	Temperature []Temperature `json:"temperature"`
}

type ReturnData struct {
	Success string `json:"success"`
	Data []TemperatureCity `json:"data"`
}

type ReturnMessage struct {
	Success string `json:"success"`
	Message string `json:"message"`
}

var Temperatures []TemperatureCity
var REGEXP_date, _ = regexp.Compile("\\d\\d.\\d\\d.\\d\\d\\d\\d")
var REGEXP_temp, _ = regexp.Compile("(-)?\\w\\s°C")

func apiTemperature(w http.ResponseWriter, r *http.Request){
	var result ReturnData
	if len(mux.Vars(r)) == 0{
		result = ReturnData{Success: "true", Data: Temperatures}
	} else {
		var result_city []TemperatureCity
		From(Temperatures).Where(func(c interface{}) bool{
			return c.(TemperatureCity).Id == mux.Vars(r)["identifier"]
		}).Select(func(c interface{}) interface{}{
			return c
		}).ToSlice(&result_city)
		result = ReturnData{Success: "true", Data: result_city}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getCityTemperature(city string){
	url := "http://www.shmu.sk/sk/?page=1&id=meteo_num_alad&mesto=" + city
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	temp := []Temperature{}
	doc.Find("h3").Each(func(_ int, item *goquery.Selection){
		date := REGEXP_date.FindString(item.Text())
		temperatures := REGEXP_temp.FindAllString(item.Next().Text(),-1)
		temp = append(temp, Temperature{Date: date, Low: temperatures[0], High: temperatures[1]})
	})
	Temperatures = append(Temperatures, TemperatureCity{Id: strings.ToLower(city), City: city, Temperature: temp})
}

func apiTemperatureScrape(w http.ResponseWriter, r *http.Request){
	Temperatures = nil
	c := colly.NewCollector()
	c.OnHTML("select", func(e *colly.HTMLElement) {
		e.ForEach("option", func(_ int, elem *colly.HTMLElement){
			getCityTemperature(elem.Text)
		})
	})
	c.Visit("http://www.shmu.sk/sk/?page=1&id=meteo_num_alad&mesto=BRATISLAVA&jazyk=sk")
	var return_message ReturnMessage
	if len(Temperatures) > 0 {
		msg := "loaded "+strconv.Itoa(len(Temperatures))+" cities"
		return_message = ReturnMessage{Success: "true", Message: msg}
	} else {
		return_message = ReturnMessage{Success: "fail", Message: "failed to load temperatures, try again"}
	}

	json.NewEncoder(w).Encode(return_message)
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/temperature", apiTemperature)
	myRouter.HandleFunc("/api/temperature/scrape", apiTemperatureScrape)
	myRouter.HandleFunc("/api/temperature/{identifier}", apiTemperature)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}


