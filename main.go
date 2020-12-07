package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type store struct {
	m map[string]string
	l sync.RWMutex
}

func serStore() *store {
	s := store{}
	s.m = make(map[string]string)
	return &s
}

func (s *store) postHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.l.Lock()
	fmt.Fprintf(w, "POST request accepted\n")
	params := mux.Vars(r) // Gets params
	// Looping through contacts and find one with the id from the params
	key := params["key"]
	value := params["value"]
	w.Write([]byte(fmt.Sprintf(`{"%s": "%s"}`, key, value)))
	if _, ok := s.m[key]; !ok {
		s.m[key] = value
		fmt.Printf("Created entry!!\n")
		fmt.Fprintf(w, "Created entry!!\n")
	} else {
		fmt.Printf("Entry exists!!\n")
		fmt.Fprintf(w, "Entry exist")
	}
	s.l.Unlock()
}
func (s *store) getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.l.RLock()
	fmt.Fprintf(w, "GET request accepted\n")
	params := mux.Vars(r) // Gets params
	key := params["key"]
	// Looping through contacts and find one with the id from the params

	if _, ok := s.m[key]; ok {
		value := s.m[key]
		fmt.Printf("Entry exists!!\n")
		w.Write([]byte(fmt.Sprintf(`{"%s": "%s"}`, key, value)))
	} else {
		fmt.Printf("Entry does not exist!!\n")
		fmt.Fprintf(w, "Value does not exist")
	}
	s.l.RUnlock()
}

func (s *store) putHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Lock()
	fmt.Fprintf(w, "PUT request accepted\n")
	params := mux.Vars(r)
	key := params["key"]
	value := params["value"]
	fmt.Fprintf(w, "Key = %s\n", key)
	if _, ok := s.m[key]; ok {
		s.m[key] = value
		fmt.Printf("Entry exists!!\n")
		fmt.Fprintf(w, "Entry exist and replaced")
	} else {
		fmt.Printf("Entry does not exist!!\n")
		fmt.Fprintf(w, "Entry does not exist")
	}
	s.l.Unlock()
}

func (s *store) delHandler(w http.ResponseWriter, r *http.Request) {
	s.l.Lock()
	fmt.Fprintf(w, "DELETE request accepted\n")
	params := mux.Vars(r)
	key := params["key"]
	fmt.Fprintf(w, "Key = %s\n", key)
	if _, ok := s.m[key]; ok {
		delete(s.m, key)
		fmt.Printf("Entry exists!!\n")
		fmt.Fprintf(w, "Deleted %s\n", key)
	} else {
		fmt.Printf("Entry does not exist!!\n")
		fmt.Fprintf(w, "Entry does not exist")
	}
	s.l.Unlock()
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func (s *store) getStored(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.m)
}

func main() {
	s := serStore()
	fmt.Println("Hello")

	r := mux.NewRouter()
	r.HandleFunc("/keys", s.getStored).Methods("GET")
	r.HandleFunc("/{key}", s.getHandler).Methods("GET")
	r.HandleFunc("/{key}/{value}", s.postHandler).Methods("POST")
	r.HandleFunc("/{key}/{value}", s.putHandler).Methods("PUT")
	r.HandleFunc("/{key}", s.delHandler).Methods("DELETE")
	r.HandleFunc("/", notFound)
	log.Fatal(http.ListenAndServe(":8080", r))
}
