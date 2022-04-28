package main

import (
	"finalProject/horoscope"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"context"
)

const (
	mongodbEndpoint = "mongodb://10.0.2.15:32488" // Find this from the Mongo container
)

var tmpl *template.Template

var userSign string 

type Post struct {
	Compatability string `bson:"compatability"`
}

type PageData struct {
	Date     string
	Sign     string
	Summary  string
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
	db := database{data: map[string]string{"name": "bday"}}
	mux := http.NewServeMux()
	mux.HandleFunc("/comp", db.compatability)
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
		userSign = newSign
		if err != nil {
			log.Fatal("Something went wrong \n")
		}
		Numbers := randomNums()
		data := PageData{
			Date:     readings.Date,
			Sign:     strings.Title(readings.Sign),
			Summary:  readings.Summary,
			LuckyNum: Numbers,
		}
		t, _ := template.ParseFiles("Website.html")
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
	numbers := randomNums()
	tempday := req.URL.Query().Get("day")
	d, _ := strconv.ParseFloat(tempday, 64)
	day := int(d)
	tempmonth := req.URL.Query().Get("month")
	m, _ := strconv.ParseFloat(tempmonth, 64)
	month := int(m)
	sign := checkBday(month, day)
	if sign == "no symbol found" {
		fmt.Fprintf(w, "Sign not found for your birthday \n")
	} else {
		readings, err := horoscope.RunCLI(sign)
		userSign = sign
		if err != nil {
			log.Fatal("Something went wrong \n")
		}
		data := PageData{
			Date:     readings.Date,
			Sign:     readings.Sign,
			Summary:  readings.Summary,
			LuckyNum: numbers,
		}
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, data)
	}
}

func (db *database) compatability(w http.ResponseWriter, req *http.Request){
	client, err := mongo.NewClient(
		options.Client().ApplyURI(mongodbEndpoint),
	)
	if err != nil {
		log.Fatal(err)
	}
	userSign = "cancer"
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	defer client.Disconnect(ctx)

	col := client.Database("myDB").Collection("signs")

	filter := bson.D{{"sign", userSign}}

	var res Post
	err = col.FindOne(ctx, filter).Decode(&res)
	//fmt.Println(res)
	fmt.Println(strings.Trim(fmt.Sprint(res), "{}"))
}

func checkBday(month, day int) string {
	if (month == 1 && day <= 19) || (month == 12 && day >= 19) {
		return "capricorn"
	}
	if (month == 1 && day >= 20) || (month == 2 && day <= 18) {
		return "aquarius"
	}
	if (month == 2 && day >= 19) || (month == 3 && day <= 20) {
		return "pisces"
	}
	if (month == 3 && day >= 21) || (month == 4 && day <= 19) {
		return "aries"
	}
	if (month == 4 && day >= 20) || (month == 5 && day <= 20) {
		return "taurus"
	}
	if (month == 5 && day >= 21) || (month == 6 && day <= 21) {
		return "gemini"
	}
	if (month == 6 && day >= 22) || (month == 7 && day <= 22) {
		return "cancer"
	}
	if (month == 7 && day >= 23) || (month == 8 && day <= 22) {
		return "leo"
	}
	if (month == 8 && day >= 23) || (month == 9 && day <= 22) {
		return "virgo"
	}
	if (month == 9 && day >= 23) || (month == 10 && day <= 23) {
		return "libra"
	}
	if (month == 10 && day >= 24) || (month == 11 && day <= 22) {
		return "scorpio"
	}
	if (month == 11 && day >= 23) || (month == 12 && day <= 20) {
		return "sagittarius"
	}
	return "no symbol found"
}
