package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

var contacts Contacts

func handler(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	fmt.Fprintf(w, "Hello %s\n", name)
}

func getContacts(w http.ResponseWriter, request *http.Request) {
	// var result map[string]interface{}
	// json.Unmarshal([]byte(byteValue), &result) // if the structure is not known
	contacts.Contacts[0].Name = "Sudhegan"
	jsonFormat, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonFormat)
}

func logging(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Logging: " + r.URL.Path)
		handler(w, r)
	}
}

func elapsedTimeForRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Println("Elapsed time for the request "+r.URL.Path+" is", time.Since(start))
		}()
		handler(w, r)
	}
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

func object(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		readObject()
		defer writeObject()
		handler(w, r)
	}
}

func compose(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, handlerFunc := range middlewares {
		handler = handlerFunc(handler)
	}
	return handler
}

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

func callCompose(handler http.HandlerFunc) http.HandlerFunc {
	return compose(handler, logging, object, elapsedTimeForRequest)
}

func generateUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return id.String()
}

func postHandler(w http.ResponseWriter, request *http.Request) {

	contactID := generateUUID()
	contactMethodID := generateUUID()
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

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/contacts", callCompose(getContacts))
	nameRouter := router.PathPrefix("/contact").Subrouter() // restricting handler under same prefix
	nameRouter.HandleFunc("/", callCompose(postHandler)).Methods("POST")
	nameRouter.HandleFunc("/{name}", callCompose(handler)).Methods("GET")
	// nameRouter.HandleFunc("/{name}", handler).Methods("GET") // restrict handler to method
	// nameRouter.HandleFunc("/{name}", handler).Host("localhost") // restrict handler to domain
	// nameRouter.HandleFunc("/{name}", handler).Schemes("http")   // restrict handler to protocol http or https
	http.ListenAndServe(":8080", router)

}