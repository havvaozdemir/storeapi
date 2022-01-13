# Store API

a RESTfull API that includes endpoints for fetching data from the provided MongoBD and getting, setting key-value pairs by in-memory store.

# How to install

git clone https://github.com/havvaozdemir/storeapi.git

# Requirements

Install the requirements listed below.

* Golang 

## How to use
* go build, run , test options are in Makefile

    >make build
    >make run
    >make test
    >make swagger
    >make build-image
    >make run-docker

* Start application with : 
    >make run
    >make run-docker

* Go to http://localhost:3000

* You can reach the services by Heroku address below:
    https://storeapi20222.herokuapp.com

## Endpoint Table

| Endpoint        | Method | Description                       |
| ----------------|--------|-----------------------------------|
| /in-memory      | POST   | Post key-value pair to memory     |
| /in-memory      | GET    | Get value of a key from memory    |
| /records        | POST   | Fetch data from Getir MongoDB     |