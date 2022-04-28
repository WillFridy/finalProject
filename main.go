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
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://10.0.2.15:32488" // Find this from the Mongo container
)

var tmpl *template.Template

var userSign string
var Numbers []int

type Post struct {
	Compatability string `bson:"compatability"`
}

type PageData struct {
	Date     string
	Sign     string
	Summary  string
	LuckyNum []int
	Images   string
	Comp     string
	Month    int
	Day      int
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
	Numbers = randomNums()
	mux := http.NewServeMux()
	mux.HandleFunc("/comp", db.compatability)
	mux.HandleFunc("/home", db.home)
	mux.HandleFunc("/read", db.read)
	mux.HandleFunc("/bday", db.bday)
	mux.HandleFunc("/about", db.about)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database struct {
	sync.Mutex
	data map[string]string
}

func (db *database) home(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()
	t, err := template.ParseFiles("Home.html")

	fmt.Println(err)
	t.Execute(w, 0)
}

func (db *database) about(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()

	t, err := template.ParseFiles("Home.html")

	fmt.Println(err)
	t.Execute(w, 0)

}
func (db *database) read(w http.ResponseWriter, req *http.Request) {
	//signList := []string{"aries", "tauras", "gemini", "cancer", "leo", "virgo", "libra",
	//	"scorpio", "sagittarius", "capricorn", "aquarius", "pisces"}
	db.Lock()
	defer db.Unlock()
	sign := userSign

	if userSign != "" {
		readings, err := horoscope.RunCLI(sign)
		userSign = sign
		if err != nil {
			log.Fatal("Something went wrong \n")
		}
		NewPicture := pictures(sign)

		data := PageData{
			Date:     readings.Date,
			Sign:     strings.Title(readings.Sign),
			Summary:  readings.Summary,
			LuckyNum: Numbers,
			Images:   NewPicture,
		}
		t, err := template.ParseFiles("AboutYou.html")
		fmt.Println(err)
		t.Execute(w, data)

	} else {
		fmt.Fprintf(w, "Sign Not Recognized \n")
	}
}

func (db *database) bday(w http.ResponseWriter, req *http.Request) {
	db.Lock()
	defer db.Unlock()

	tempday := req.URL.Query().Get("day")
	d, _ := strconv.ParseFloat(tempday, 64)
	day := int(d)

	tempmonth := req.URL.Query().Get("month")
	m, _ := strconv.ParseFloat(tempmonth, 64)
	month := int(m)

	if userSign == "" {
		sign := checkBday(month, day)
		readings, err := horoscope.RunCLI(sign)
		NewPicture := pictures(sign)
		userSign = sign
		if err != nil {
			log.Fatal("Something went wrong \n")
		}
		data := PageData{
			Date:     readings.Date,
			Sign:     strings.Title(readings.Sign),
			Summary:  readings.Summary,
			LuckyNum: Numbers,
			Images:   NewPicture,
		}
		t, _ := template.ParseFiles("BDay.html")
		t.Execute(w, data)
	} else {
		readings, _ := horoscope.RunCLI(userSign)
		NewPicture := pictures(userSign)
		data := PageData{
			Date:     readings.Date,
			Sign:     strings.Title(readings.Sign),
			Summary:  readings.Summary,
			LuckyNum: Numbers,
			Images:   NewPicture,
		}
		t, _ := template.ParseFiles("BDay.html")
		t.Execute(w, data)
	}
}

func (db *database) compatability(w http.ResponseWriter, req *http.Request) {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(mongodbEndpoint),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	col := client.Database("myDB").Collection("signs")

	filter := bson.D{{"sign", userSign}}

	var res Post
	err = col.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		log.Fatal(err)
	}

	output := strings.Trim(fmt.Sprint(res), "{}")
	data := PageData{
		Comp: output,
	}
	t, _ := template.ParseFiles("Website.html")
	t.Execute(w, data)

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

func pictures(sign string) string {
	if sign == "capricorn" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-capricorn-265x300.png"
	}
	if sign == "aquarius" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-aquarius-300x145.png"
	}
	if sign == "pisces" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-pisces-150x150.png"
	}
	if sign == "aries" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-aries-150x150.png"
	}
	if sign == "taurus" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-taurus-285x300.png"
	}
	if sign == "gemini" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-gemini-255x300.png"
	}
	if sign == "cancer" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-cancer-300x207.png"
	}
	if sign == "leo" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-leo-199x300.png"
	}
	if sign == "virgo" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-virgo-279x300.png"
	}
	if sign == "libra" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-libra-300x176.png"
	}
	if sign == "scorpio" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-scorpio-262x300.png"
	}
	if sign == "sagittarius" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-sagittarius-300x232.png"
	}
	return "Not Found"
}

func description(sign string) string {
	if sign == "capricorn" {
		return "Capricorns are masters of discipline. The wringing of the hands, the constant reminders, the exacting structure, the ever-increasing goals, the tidal wave of self-criticism that lasts forever. /n They are the ultimate perfectionist. They can be so absorbed in their own internal monologue that it becomes impossible to get them to look away from themselves. Capricorns are often called “workaholics.” "
	}
	if sign == "aquarius" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-aquarius-300x145.png"
	}
	if sign == "pisces" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-pisces-150x150.png"
	}
	if sign == "aries" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-aries-150x150.png"
	}
	if sign == "taurus" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-taurus-285x300.png"
	}
	if sign == "gemini" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-gemini-255x300.png"
	}
	if sign == "cancer" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-cancer-300x207.png"
	}
	if sign == "leo" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-leo-199x300.png"
	}
	if sign == "virgo" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-virgo-279x300.png"
	}
	if sign == "libra" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-libra-300x176.png"
	}
	if sign == "scorpio" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-scorpio-262x300.png"
	}
	if sign == "sagittarius" {
		return "https://astrostyle.com/wp-content/uploads/2020/07/signs-glyphs-sagittarius-300x232.png"
	}
	return "Not Found"
}
