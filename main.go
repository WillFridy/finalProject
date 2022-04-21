package main

import (
	"finalProject/horoscope"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	db := database{data: map[string]string{"libra": "ligma"}}
	mux := http.NewServeMux()
	mux.HandleFunc("/create", db.create)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database struct {
	sync.Mutex
	data map[string]string
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	newSign := req.URL.Query().Get("sign")
	readings, err := horoscope.RunCLI(newSign)
	if err != nil {
		log.Fatal("shit went wrong")
	}
	fmt.Fprintf(w, "%s \n", readings)
}
