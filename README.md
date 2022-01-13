# Store API

a RESTfull API that includes endpoints for fetching data from the provided MongoBD and getting, setting key-value pairs by in-memory store.

# How to install

git clone [github address]

# Requirements

Install the requirements listed below.

* Golang 

## How to use
* go build, run , test options are in Makefile
    >Make build
    >Make run
    >Make test
    >Make swagger
    >Make build-image
    >Make run-docker

* Start application with : 
    make run

* Go to http://localhost:8080

* You can reach the services by Heroku address below:
    [heroku address]

## Endpoint Table

| Endpoint        | Method | Description                       |
| ----------------|--------|-----------------------------------|
| /in-memory      | POST   | Post key-value pair to memory     |
| /in-memory      | GET    | Get value of a key from memory    |
| /records        | POST   | Fetch data from Getir MongoDB     |