package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"utilities"

	"github.com/gorilla/mux"
)

type Contacts struct {
	Contacts []Contact `json:"contacts"`
}

type Contact struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Contactmethods []ContactMethod `json:"contactMethods"`
}

type ContactMethod struct {
	ID         string `json:"id"`
	MethodType string `json:"methodType"`
	Value      string `json:"value"`
}

var contacts Contacts

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
	vars := mux.Vars(request)
	name := vars["name"]
	fmt.Fprintf(w, "Hello %s\n", name)
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
	contactMethod1 := ContactMethod{
		contactMethodID,
		"email",
		"jeevasudhegan1198@gmail.com",
	}
	contact := Contact{
		contactID,
		"Jeeva",
		[]ContactMethod{contactMethod1},
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

func main() {

	readObject()
	// defer writeObject()
	router := mux.NewRouter()
	router.HandleFunc("/contacts", utilities.CallCompose(getContacts))
	router.HandleFunc("/logout", closeHandler)
	nameRouter := router.PathPrefix("/contact").Subrouter() // restricting handler under same prefix
	nameRouter.HandleFunc("/", utilities.CallCompose(postHandler)).Methods("POST")
	nameRouter.HandleFunc("/{name}", utilities.CallCompose(handler)).Methods("GET")
	// nameRouter.HandleFunc("/{name}", handler).Methods("GET") // restrict handler to method
	// nameRouter.HandleFunc("/{name}", handler).Host("localhost") // restrict handler to domain
	// nameRouter.HandleFunc("/{name}", handler).Schemes("http")   // restrict handler to protocol http or https
	http.ListenAndServe(":8080", router)

}
