package main

import (
	"finalProject/horoscope"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

var tmpl *template.Template

type SignInfo struct {
	Data     string
	Sign     string
	Summery  string
	LuckyNum []int
}

type PageData struct {
	Title    string
	Date     string
	Sign     string
	Summery  string
	LuckyNum []int
}

func randomNums() (nums []int) {
	i := 0
	for i < 5 {
		randNum := rand.Intn(100)
		nums = append(nums, randNum)
		i++
	}
	return nums
}

func checkSign(signList []string, str string) bool {
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
	mux.HandleFunc("/read", db.read)
	mux.HandleFunc("/bday", db.bday)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database struct {
	sync.Mutex
	data map[string]string
}

func (db *database) read(w http.ResponseWriter, req *http.Request) {
	signList := []string{"aries", "tauras", "gemini", "cancer", "leo", "virgo", "libra",
		"scorpio", "sagittarius", "capricorn", "aquarius", "pisces"}
	db.Lock()
	defer db.Unlock()
	newSign := req.URL.Query().Get("sign")

	if checkSign(signList, newSign) {
		readings, err := horoscope.RunCLI(newSign)
		if err != nil {
			log.Fatal("Something went horribly wrong \n")
		}
		Numbers := randomNums()
		data := PageData{
			Title:    readings.Sign,
			Date:     readings.Date,
			Sign:     readings.Sign,
			Summery:  readings.Summary,
			LuckyNum: Numbers,
		}
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, data)

		// fmt.Fprintf(w, "Date: %s \n", readings.Date)
		// fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		// fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		// fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
	} else {
		fmt.Fprintf(w, "Sign Not Recognized \n")
	}
}

func (db *database) bday(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	Numbers := randomNums()
	tempday := req.URL.Query().Get("day")
	d, _ := strconv.ParseFloat(tempday, 64)
	day := int(d)
	tempmonth := req.URL.Query().Get("month")
	m, _ := strconv.ParseFloat(tempmonth, 64)
	month := int(m)

	if (month == 1 && day <= 19) || (month == 12 && day >= 19) {
		readings, _ := horoscope.RunCLI("capricorn")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 1 && day >= 20) || (month == 2 && day <= 18) {
		readings, _ := horoscope.RunCLI("aquarius")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 2 && day >= 19) || (month == 3 && day <= 20) {
		readings, _ := horoscope.RunCLI("pisces")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 3 && day >= 21) || (month == 4 && day <= 19) {
		readings, _ := horoscope.RunCLI("aries")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 4 && day >= 20) || (month == 5 && day <= 20) {
		readings, _ := horoscope.RunCLI("taurus")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 5 && day >= 21) || (month == 6 && day <= 21) {
		readings, _ := horoscope.RunCLI("gemini")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 6 && day >= 22) || (month == 7 && day <= 22) {
		readings, _ := horoscope.RunCLI("cancer")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 7 && day >= 23) || (month == 8 && day <= 22) {
		readings, _ := horoscope.RunCLI("leo")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 8 && day >= 23) || (month == 9 && day <= 22) {
		readings, _ := horoscope.RunCLI("virgo")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 9 && day >= 23) || (month == 10 && day <= 23) {
		readings, _ := horoscope.RunCLI("libra")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 10 && day >= 24) || (month == 11 && day <= 22) {
		readings, _ := horoscope.RunCLI("scorpio")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
	if (month == 11 && day >= 23) || (month == 12 && day <= 20) {
		readings, _ := horoscope.RunCLI("sagitarrius")
		fmt.Fprintf(w, "Date: %s \n", readings.Date)
		fmt.Fprintf(w, "Sign: %s \n", readings.Sign)
		fmt.Fprintf(w, "Summary: %s \n", readings.Summary)
		fmt.Fprintf(w, "Lucky Numbers: %d \n", Numbers)
		fmt.Fprintf(w, "Birthday:  %d/%d \n", month, day)
	}
}
