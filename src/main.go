package main

import (
	"dto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"utilities"

	"github.com/gorilla/websocket"

	"github.com/gorilla/mux"
)

var contacts dto.Contacts

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func readObject() {
	log.Println("Creating the object from file...")
	contact, err := os.Open("contacts.json")
	if err != nil {
		log.Fatal(err)
	}
	defer contact.Close()
	byteValue, _ := ioutil.ReadAll(contact)
	json.Unmarshal(byteValue, &contacts)
}

func writeObject() {
	log.Println("Writing the object from file...")
	jsonFormat, err := json.MarshalIndent(contacts, "", "  ") // indented json
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("contacts.json", jsonFormat, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// func getObject(handler http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		readObject()
// 		defer writeObject()
// 		handler(w, r)
// 	}
// }

// Handler functions
func handler(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request) // this returns a map
	if len(vars) == 0 {
		data := dto.TodoPageData{
			PageTitle: "My TODO list",
			Todos: []dto.Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
			},
		}
		tmpl := template.Must(template.ParseFiles("layout.html"))
		tmpl.Execute(w, data)
	} else {
		name := vars["name"]
		fmt.Fprintf(w, "Hello %s\n", name)
	}
}

func getContacts(w http.ResponseWriter, request *http.Request) {
	// var result map[string]interface{}
	// json.Unmarshal([]byte(byteValue), &result) // if the structure is not known
	if len(contacts.Contacts) != 0 {
		contacts.Contacts[0].Name = "Sudhegan"
	}
	jsonFormat, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonFormat)
}

func postHandler(w http.ResponseWriter, request *http.Request) {

	contactID := utilities.GenerateUUID()
	contactMethodID := utilities.GenerateUUID()
	contactMethod1 := dto.ContactMethod{
		contactMethodID,
		"email",
		"jeevasudhegan1198@gmail.com",
	}
	contact := dto.Contact{
		contactID,
		"Jeeva",
		[]dto.ContactMethod{contactMethod1},
	}
	contacts.Contacts = append(contacts.Contacts, contact)
	jsonFormat, err := json.MarshalIndent(contact, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonFormat)
}

func closeHandler(w http.ResponseWriter, r *http.Request) {
	writeObject()
	log.Println("Stopping the server...")
	os.Exit(0)
}

func chatsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	if err != nil {
		log.Fatal("unable to connect to websocket")
		return
	}

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

func main() {

	readObject()
	// defer writeObject()
	fs := http.FileServer(http.Dir("/assets/"))
	router := mux.NewRouter()
	http.HandleFunc("/socket", chatsocket)
	router.Handle("/static/", http.StripPrefix("/static/", fs))
	router.HandleFunc("/echo", chatsocket)

	router.HandleFunc("/websockets", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})
	router.HandleFunc("/contacts", utilities.CallCompose(getContacts))
	router.HandleFunc("/logout", closeHandler)
	nameRouter := router.PathPrefix("/contact").Subrouter() // restricting handler under same prefix
	nameRouter.HandleFunc("/", utilities.CallCompose(postHandler)).Methods("POST")
	nameRouter.HandleFunc("/", utilities.CallCompose(handler)).Methods("GET")
	nameRouter.HandleFunc("/{name}", utilities.CallCompose(handler)).Methods("GET")
	// nameRouter.HandleFunc("/{name}", handler).Methods("GET") // restrict handler to method
	// nameRouter.HandleFunc("/{name}", handler).Host("localhost") // restrict handler to domain
	// nameRouter.HandleFunc("/{name}", handler).Schemes("http")   // restrict handler to protocol http or https
	http.ListenAndServe(":8080", router)

}
