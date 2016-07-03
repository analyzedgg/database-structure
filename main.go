package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	COUCHDB_URL = "http://localhost:5984/"
	OK = "OK"
	ALREADYEXIST = "Already exists"
	UNKNOWNERROR = "Unknown error occurred(%s)"
)

func createDb(name string, ch chan string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			os.Exit(1)
		}
	}()

	request, err := http.NewRequest(http.MethodPut, COUCHDB_URL + name, nil)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		message := fmt.Sprintf("Error occurred, is the CouchDB server running?")
		panic(message)
	}

	switch response.StatusCode {
	case http.StatusCreated:
		ch <- OK
	case http.StatusPreconditionFailed:
		ch <- ALREADYEXIST
	default:
		ch <- fmt.Sprintf(UNKNOWNERROR, response.StatusCode)
	}
}

func createStructure() {
	summonerChan := make(chan string)
	matchesChan := make(chan string)

	go createDb("summoner-db", summonerChan)
	go createDb("matches-db", matchesChan)

	fmt.Println("SummonerDb creation: ", <-summonerChan)
	fmt.Println("MatchesDb creation: ", <-matchesChan)

	fmt.Println("Done")
}

func main() {
	fmt.Printf("About to create database structure on CouchDB: %s\n", COUCHDB_URL)

	createStructure()
}