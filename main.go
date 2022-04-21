package main

import (
	"finalProject/horoscope"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

func randomNums() (nums []int) {
	i := 0
	for i < 5 {
		randNum := rand.Intn(100)
		nums = append(nums, randNum)
		i++
	}
	return nums
}

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
	Numbers := randomNums()
	fmt.Fprintf(w, "Date: %s \n", readings.Date)
	fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
	fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
	fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
}
