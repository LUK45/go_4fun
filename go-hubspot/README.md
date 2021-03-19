# Hubspot integration task 
Syncs contacts and tickets from hubspot to local mysql db. 

* For contacts it stores: id, first name, last name, email, phone
* For ticket: id, content, hubspot_owner_id  (Hubspot API doesn't support method to get ticket owners name)

## Usage

* create local mysql db `hubspot_db`
* get hubspot api key and set hubspot api key and db credentials in `main.go` lines 17-19
* download go deps
* `go build` 
* `./go-hubspot`


* This starts the program. It starts API on `localhost:10000` which can be used to get data from db. 


* It creates tables in db if they don't exist. 
* Right after the start it runs first sync of contacts and tickets to local db.
* Next it syncs data periodically every 60 seconds. This period can be changed on `line 25 in main.go `

* Data stored in db can be read as json through api:

```
$ curl localhost:10000/api/tickets | jq
[
  {
    "id": "327186038",
    "properties": {
      "content": "jkhkjfhdvkjhdfkjvhdfkjvhdfkj",
      "hubspot_owner_id": "64284690"
    }
  },
  {
    "id": "329106860",
    "properties": {
      "content": "lololocdscsdcsdcsdcsd",
      "hubspot_owner_id": "64284690"
    }
  }
]

$ curl localhost:10000/api/contacts | jq
[
  {
    "id": "1",
    "properties": {
      "firstname": "Maria",
      "lastname": "Johnson (Sample Contact)",
      "email": "emailmaria@hubspot.com",
      "phone": "12312312"
    }
  },
  {
    "id": "51",
    "properties": {
      "firstname": "Brian",
      "lastname": "Halligan (Sample Contact)",
      "email": "bh@hubspot.com",
      "phone": "32423432"
    }
  },
  {
    "id": "101",
    "properties": {
      "firstname": "lukas",
      "lastname": "masar",
      "email": "11111.m1frfrfrasar@gmail.com",
      "phone": "0904111729"
    }
  }
]
```

There is an option of manual start of sync: 

```
$ curl localhost:10000/api/hubspot/contacts | jq
{
  "success": "success",
  "message": "synced 3 records"
}

$ curl localhost:10000/api/hubspot/tickets | jq
{
  "success": "success",
  "message": "synced 2 records"
}
```
