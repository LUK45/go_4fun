package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
	_ "github.com/go-sql-driver/mysql"

)

var DB_USER = "DB_USER"
var DB_PASSWORD = "DB_PASSWORD"
var HUBSPOT_API_KEY = "HUBSPOT_API_KEY_GOES_HERE"
var HUBSPOT_API_URL = "https://api.hubapi.com/crm/v3/objects/"
var CONTACTS_HUBSPOT_LIMIT = "100"
var TICKETS_HUBSPOT_LIMIT = "100"
var CONTACTS_URL = HUBSPOT_API_URL+"contacts?limit="+CONTACTS_HUBSPOT_LIMIT+"&properties=firstname,lastname,email,phone&hapikey="+HUBSPOT_API_KEY
var TICKETS_URL = HUBSPOT_API_URL+"tickets?limit="+TICKETS_HUBSPOT_LIMIT+"&properties=content,hubspot_owner_id&hapikey="+HUBSPOT_API_KEY
var SYNC_INTERAVAL_SECONDS = uint64(60)
var DB *sql.DB

type ReturnMessage struct {
	Success string `json:"success"`
	Message string `json:"message"`
}

type ContactProperties struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type TicketProperties struct {
	Content string `json:"content"`
	HubspotOwnerId string `json:"hubspot_owner_id"`
}

type Ticket struct {
	Id string `json:"id"`
	Properties TicketProperties `json:"properties"`
}

type Contact struct {
	Id string `json:"id"`
	Properties ContactProperties `json:"properties"`
}

type HPNext struct {
	Link string `json:"link"`
}

type HubspotPaging struct {
	Next HPNext `json:"next"`
}

type ContactHubspotRespone struct {
	Results []Contact `json:"results"`
	Paging *HubspotPaging `json:"paging"`
}

type TicketHubspotRespone struct {
	Results []Ticket `json:"results"`
	Paging *HubspotPaging `json:"paging"`
}

func syncFromHubspot(url string, entity string) int {
	next := true
	count := 0
	for next {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		data, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		switch entity {
			case "contact":
				var contacts ContactHubspotRespone
				err2 := json.Unmarshal(data, &contacts)
				if err2 != nil {
					fmt.Println(err2)
				}
				for _, c := range contacts.Results {
					err3 := store.CreateContact(&c)
					if err3 != nil {
						fmt.Println(err3)
					}
					count ++
				}
				if contacts.Paging != nil {
					url = contacts.Paging.Next.Link + "&hapikey=" + HUBSPOT_API_KEY
				} else {
					next = false
				}
			case "ticket":
				var tickets TicketHubspotRespone
				err2 := json.Unmarshal(data, &tickets)
				if err2 != nil {
					fmt.Println(err2)
				}
				for _, t := range tickets.Results {
					err3 := store.CreateTicket(&t)
					if err3 != nil {
						fmt.Println(err3)
					}
					count ++
				}
				if tickets.Paging != nil {
					url = tickets.Paging.Next.Link + "&hapikey=" + HUBSPOT_API_KEY
				} else {
					next = false
				}
			}
	}
	return count
}

func apiHubspot(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var count int
	entity :=  mux.Vars(r)["entity"];
	switch entity {
		case "contacts":
			count = syncFromHubspot(CONTACTS_URL,"contact")
		case "tickets":
			count = syncFromHubspot(TICKETS_URL,"ticket")
	}
	return_message := ReturnMessage{Success: "success", Message: "synced "+strconv.Itoa(count)+" records"}
	json.NewEncoder(w).Encode(return_message)

}

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	entity :=  mux.Vars(r)["entity"];
	switch entity {
		case "contacts":
			result, err := store.GetContacts()
			if err != nil {
				fmt.Println(fmt.Errorf("Error: %v", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(result)

		case "tickets":
			result, err := store.GetTickets()
			if err != nil {
				fmt.Println(fmt.Errorf("Error: %v", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(result)
	}

}

func executeCronJob() {
	gocron.Every(SYNC_INTERAVAL_SECONDS).Second().Do(syncFromHubspot, CONTACTS_URL, "contact")
	gocron.Every(SYNC_INTERAVAL_SECONDS).Second().Do(syncFromHubspot, TICKETS_URL, "ticket")
	<- gocron.Start()
}

func main() {
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(127.0.0.1:3306)/hubspot_db")
	if err != nil {
		panic(err.Error())
	}
	InitStore(&dbStore{db: db})
	store.CreateTables()
	go syncFromHubspot(CONTACTS_URL,"contact")
	go syncFromHubspot(TICKETS_URL,"ticket")
	go executeCronJob()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/hubspot/{entity}", apiHubspot)
	myRouter.HandleFunc("/api/{entity}", getHandler)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
	defer db.Close()
}
