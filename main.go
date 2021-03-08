package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "github.com/gorilla/mux"
	"os"
	"time"
)

type Message struct {
	Text        string    `json:"text"`
	Dt          time.Time `json:"dt"`
	Source      User      `json:"source"`
	Destination User      `json:"destination"`
}

type User struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
}

var Messages []Message

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage of Mirk!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllMessages")
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(Messages)
	case "POST":
		fmt.Fprintf("POST method not supported yet")
	default:
		fmt.Fprintf("Only GET and POST methods are supported")
	}
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/messages", returnAllMessages)
	log.Fatal(http.ListenAndServe(GetPort(), nil))
}

func main() {
	alice := User{Id: 0, Nickname: "Alice"}
	bob := User{Id: 1, Nickname: "Bob"}
	now := time.Now()

	Messages = []Message{
		Message{Text: "Hello world", Dt: now, Source: alice, Destination: bob},
	}
	handleRequests()
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
