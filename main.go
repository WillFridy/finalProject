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

func checkSign(signList []string,str string) bool {
	for _, v := range signList {
		if v == str {
			return true
		}
	}
	return false
}
func main() {
	db := database{data: map[string]string{"fill": "fill"}}
	mux := http.NewServeMux()
	mux.HandleFunc("/create", db.create)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database struct {
	sync.Mutex
	data map[string]string
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	signList := []string{"aries", "aauras", "gemini", "cancer", "leo", "virgo", "libra", "scorpio", "sagittarius", "capricorn", "aquarius", "picses"}
	db.Lock()
	defer db.Unlock()
	newSign := req.URL.Query().Get("sign")

	if checkSign(signList, newSign) {
		readings, err := horoscope.RunCLI(newSign)
		if err != nil {
			log.Fatal("Something went horribly wrong \n")
		}
		Numbers := randomNums()
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
	} else {
		log.Fatal("Sign Not Recognized \n")
	}
	
}
