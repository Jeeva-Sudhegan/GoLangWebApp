package dto

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
