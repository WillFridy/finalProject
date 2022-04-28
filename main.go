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
	Des      string
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
		info := description(userSign)
		data := PageData{
			Date:     readings.Date,
			Sign:     strings.Title(readings.Sign),
			LuckyNum: Numbers,
			Images:   NewPicture,
			Des:      info,
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
		return "Capricorns are masters of discipline. The wringing of the hands, the constant reminders, the exacting structure, the ever-increasing goals, the tidal wave of self-criticism that lasts forever. \n They are the ultimate perfectionist. They can be so absorbed in their own internal monologue that it becomes impossible to get them to look away from themselves. Capricorns are often called “workaholics.” "
	}
	if sign == "aquarius" {
		return "Aquarians are archetypical outcasts. This doesn't mean they're loners. In fact, they thrive in large groups—charming you with their peculiar senses of humor, intriguing you with fun facts about the history of disposable straws, or convincing you to join their reading group. The alienation they feel is often self-imposed—a result of their knee-jerk contrarianism, rather than a lack of social intelligence. They try to be weird. They hang grapefruit rinds from the wall and call it art, they pretend to actually like noise music, they saturate their internal monologues with SAT words."
	}
	if sign == "pisces" {
		return "If you are asking this question, you're probably a Pisces who is insecure. If you don't feel valued, it's not because being a Pisces is bad, but because society as a whole generally undervalues “soft” skills like intuition and sensitivity. Your challenge is to start viewing these things as talents instead of impediments."
	}
	if sign == "aries" {
		return "At their core, Aries do what they want and do things their way. They are unafraid of conflict, highly competitive, honest and direct. An Aries is not weighed down by the freedom of choice, and is perhaps the sign that is least conflicted about what they want. They throw themselves at the world eagerly and without fear. It is one of their most commendable qualities, but also what causes them a great deal of pain and grief."
	}
	if sign == "taurus" {
		return "Taureans are the human equivalent of moss. A handmade wooden chair. They are normally satisfied with the way things are. They embody stability. Sitting in a patch of grass admiring the breeze. When everything else seems to be falling apart, Taureans are an oasis of calm, a rock of dependability. Practical knowledge and experience is their modus operandi."
	}
	if sign == "gemini" {
		return "Geminis are very intelligent and pick up knowledge quickly. They are perceptive, analytical, and often very funny. They have an unreserved and childlike curiosity, always asking new questions. Geminis have an uncanny ability to size up a person's character in a matter of seconds, even if they only just met them. If someone's bluffing, they'll be the first to notice. They are great communicators, very responsive and sensitive listeners."
	}
	if sign == "cancer" {
		return "A Cancer's personality is like wading chest deep in a lake of warm water. It feels sparkling and cool while it's touching the body, but you know that if you were to dive in, it would feel warm. Cancer's self-awareness is like the tides. They're constantly moving in and out of focus. Their personality is layered. They have many moods, some of which are contradictory, but they also have a deep, core self that persists."
	}
	if sign == "leo" {
		return "Leos are bold, warm, and loving. They are also the ultimate showmen. They can dazzle with the theatrical flair of a Broadway star and the charisma of a politician. They are captivating personalities. They have a way with words, and can speak eloquently on just about any topic, no matter how quickly they've just been introduced to it."
	}
	if sign == "virgo" {
		return "It's true that Virgos are very particular, but that doesn't necessarily mean that they keep neat spaces. Their particularities and habits don't necessarily line up with traditional views of cleanliness. They could live in what looks like a Tasmanian devil-style dust storm ruin, but still impose a “no shoes in the house” or “no outside clothes on the bed” rule. Maybe their house looks cluttered, but they still know where everything is. Everything has its place. Virgos prefer to exist in organized spaces, but put their service orientation over their own comfort. This can mean that a Virgo is too busy fixing the lives of those around them to put much work into providing for their own needs. They're rarely motivated by their own self-interest."
	}
	if sign == "libra" {
		return "Libras are difficult to really understand because they seem so contradictory on the surface. They're simultaneously extroverted and introverted, strategic and spontaneous, focused and intuitive. This variability makes it difficult to pin down their true character. They are an entire constellation of personalities. Libras are different depending on who they're around."
	}
	if sign == "scorpio" {
		return "The Scorpio personality is a profound chasm of infinite complexity (or at least how they project themselves). They are difficult people to get to know. They are psychological trap doors. They socialize from behind a double-sided mirror, always scanning, reading you while you can only see your own reflection. They prefer to be the people asking the questions. They remove your skin with their perceptive scalpel and take inventory of your pulsing viscera. They probe and push. They know the little things that make you tick. Your pressure points. The subtle ways to procure the answer they're seeking. They are keenly aware of power, its flows, and their position within its matrix."
	}
	if sign == "sagittarius" {
		return "Sagittarius is the ultimate empiricist. They will always choose principles over feelings and will often question who they are. They move from job to job, philosophy to philosophy, belief to belief. They are explorers of the human condition and are unafraid of change. Sagittarians feel like the world is their playground. They love to explore the unknown. They want to understand how the world works."
	}
	return "Not Found"
}
