package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	// "strings"
)

type Message struct {
	Text        string    `json:"txt"`
	Dt          time.Time `json:"dt"`
	Source      User      `json:"source"`
	Destination User      `json:"destination"`
	ChatId      int       `json:"chat"`
}

type User struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
}

type Chat struct {
	Id    int    `json:"id"`
	Users []User `json:"users"`
}

// dummy DB
var Messages []Message
var Chats []Chat
var Users []User

func getLastUserId() int {
	if len(Users) == 0 {
		return 0
	}
	return Users[len(Users)-1].Id + 1
}

func validateUser(user User, password string) bool {
	for _, u := range Users {
		if u.Id == user.Id && u.Nickname == user.Nickname && len(password) != 0 {
			return true
		}
	}
	return false
}

func getLastChatId() int {
	if len(Chats) == 0 {
		return 0
	}
	return Chats[len(Chats)-1].Id + 1
}

//TODO add errors
func chatsWithUser(user User) []Chat {
    var chats []Chat
    for _, chat := range Chats {
        for _, chatUser := range chat.Users {
            if chatUser.Id == user.Id {
                chats = append(chats, chat)
            }
        }
    }
    return chats
}

//TODO add errors
func findUserById(id int) User {
    for _, u := range Users {
        if u.Id == id {
            return u
        }
    }
    return User{}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage of Mirk!")
	fmt.Println("Endpoint Hit: homePage")
}

func messagesRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: messagesRoute")
	switch r.Method {
	//getChatMessages(chatId)
	case "GET":

		json.NewEncoder(w).Encode(Messages)

		//sendMessage(m)
	case "POST":
		fmt.Fprintf(w, "sourceId: %s\n", r.FormValue("srcId"))
		fmt.Fprintf(w, "destinationId: %s\n", r.FormValue("destId"))

	//deleteMessage(m)
	case "DELETE":

	default:
		fmt.Fprintf(w, "Only GET and POST methods are supported")
	}
}

func usersRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: usersRoute")
	switch r.Method {
	//get all users
	case "GET":
		json.NewEncoder(w).Encode(Users)
	//create new user
	case "POST":
		Users = append(Users, User{Id: getLastUserId(), Nickname: r.FormValue("nickname")})
		fmt.Fprintf(w, "Added new user with nickname: %s\n", r.FormValue("nickname"))
	//delete user
	case "DELETE":

	default:
		fmt.Fprintf(w, "Only GET and POST methods are supported")
	}
}

func usersValidateRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: usersValidateRoute")
	switch r.Method {
	//validate user
	case "GET":
		userId, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			fmt.Println("Err while parsing user id")
		}
		user := User{Nickname: r.FormValue("nickname"), Id: userId}
		if validateUser(user, r.FormValue("password")) {
			fmt.Fprintf(w, "Valid user")
		} else {
			http.Error(w, "Bad user", http.StatusForbidden)
		}
	default:
		fmt.Fprintf(w, "Only GET methods are supported")
	}
}

func chatsRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: chatsRoute")
	switch r.Method {
	//get all chats for admin and for user
	case "GET":
        // get all chats for user
        if len(r.FormValue("userId"))!=0 {
            userId, err := strconv.Atoi(r.FormValue("userId"))
            if err != nil {
                fmt.Println("Err while parsing user id")
                fmt.Fprintf(w, "Please specify some correct user id")
            }
            json.NewEncoder(w).Encode(chatsWithUser(findUserById(userId)))    
        } else {
            if len(r.FormValue("adminKey"))!=0 {
                json.NewEncoder(w).Encode(Chats)
            }
            fmt.Fprintf(w, "Please specify some userId")
        }
	//create new chat
	case "POST":
		str := r.FormValue("users")
		
        var users []User
        json.Unmarshal([]byte(str), &users)

		chat := Chat{Id: getLastChatId(), Users: users}
		Chats = append(Chats, chat)
		
        fmt.Fprintf(w, "newChat:\n id: %v\n users:\n", chat.Id)
		for _, u := range users {
			fmt.Fprintf(w, "  id: %v\n  nickname: %s\n", u.Id, u.Nickname)
		}
	//delete chat
	case "DELETE":

	default:
		fmt.Fprintf(w, "Only GET and POST methods are supported")
	}
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/api/v1/messages", messagesRoute)
	http.HandleFunc("/api/v1/users", usersRoute)
	http.HandleFunc("/api/v1/users/validate", usersValidateRoute)
	http.HandleFunc("/api/v1/chats", chatsRoute)
	log.Fatal(http.ListenAndServe(GetPort(), nil))
}

func main() {
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
