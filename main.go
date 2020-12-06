package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type store struct {
	m map[string]string
	l sync.RWMutex
}

func serStore() *store {
	return &store{}
}

func (s *store) postHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/write" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		s.l.Lock()
		fmt.Fprintf(w, "POST request accepted\n")
		key := r.FormValue("Key")
		value := r.FormValue("Value")
		fmt.Fprintf(w, "Key = %s\n", key)
		fmt.Fprintf(w, "Value = %s\n", value)
		if _, ok := s.m[key]; !ok {
			s.m[key] = value
			fmt.Printf("Created entry!!\n")
			w.Write([]byte(`{"message": "Created"}`))
		} else {
			fmt.Printf("Entry exists!!\n")
			w.Write([]byte(`{"message": "Already exists"}`))
		}
		s.l.Unlock()
	case "GET":
		s.l.RLock()
		fmt.Fprintf(w, "POST request accepted\n")
		key := r.URL.Query().Get("key")
		fmt.Fprintf(w, "Key = %s\n", key)
		if _, ok := s.m[key]; ok {
			value := s.m[key]
			fmt.Printf("Entry exists!!\n")
			fmt.Fprintf(w, "Value = %s\n", value)
			w.Write([]byte(`{"message": "Exists"}`))
		} else {
			fmt.Printf("Entry does not exist!!\n")
			fmt.Fprintf(w, "Value does not exist")
			w.Write([]byte(`{"message": "Does not exist"}`))
		}
		s.l.RUnlock()
	case "PUT":
		s.l.Lock()
		fmt.Fprintf(w, "POST request accepted\n")
		key := r.FormValue("Key")
		value := r.FormValue("Value")
		fmt.Fprintf(w, "Key = %s\n", key)
		if _, ok := s.m[key]; ok {
			s.m[key] = value
			fmt.Printf("Entry exists!!\n")
			fmt.Fprintf(w, "Value = %s\n", value)
			w.Write([]byte(`{"message": "Exists"}`))
		} else {
			fmt.Printf("Entry does not exist!!\n")
			fmt.Fprintf(w, "Entry does not exist")
			w.Write([]byte(`{"message": "Does not exist"}`))
		}
		s.l.Unlock()
	case "DELETE":
		s.l.Lock()
		fmt.Fprintf(w, "POST request accepted\n")
		key := r.URL.Path[7:]
		fmt.Fprintf(w, "Key = %s\n", key)
		if _, ok := s.m[key]; ok {
			delete(s.m, key)
			fmt.Printf("Entry exists!!\n")
			fmt.Fprintf(w, "Deleted %s\n", key)
			w.Write([]byte(`{"message": "Deleted"}`))
		} else {
			fmt.Printf("Entry does not exist!!\n")
			fmt.Fprintf(w, "Entry does not exist")
			w.Write([]byte(`{"message": "Does not exist"}`))
		}
		s.l.Unlock()
	default:
		fmt.Fprintf(w, "Sorry; POST, GET, PUT, DELETE methods are supported for this url.")
	}
}

func main() {
	s := serStore()

	http.HandleFunc("/", s.postHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
