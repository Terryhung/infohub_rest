package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/Terryhung/infohub_rest/infohub_user"
	"github.com/Terryhung/infohub_rest/mongo_lib"
	"github.com/Terryhung/infohub_rest/news"
	"github.com/Terryhung/infohub_rest/redis_lib"
	"github.com/Terryhung/infohub_rest/stock"
	"github.com/Terryhung/infohub_rest/user_event"
	"github.com/ant0ine/go-json-rest/rest"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

// Structure for DB login and API Respond
type User struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
}

type Account struct {
	Mongo_users []User
}

type Result struct {
	Message      string      `json:"message"`
	News_result  interface{} `json:"news_results"`
	Video_result interface{} `json:"video_results"`
	Image_result interface{} `json:"image_results"`
}

type Respond struct {
	Code   int
	Result Result
}

// Routing
func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	api.Use(&rest.GzipMiddleware{})
	router, err := rest.MakeRouter(
		rest.Get("/get_news", GetNews),
		rest.Get("/get_video", GetVideo),
		rest.Get("/get_image", GetImage),
		rest.Get("/get_all", GetAll),
		rest.Get("/ping", Ping),
		rest.Post("/v1/user_event", PostUserEvent),
		rest.Get("/v1/keyword", GetNewsByKeyword),
		rest.Get("/v1/stocks", GetStockList),
	)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8787", api.MakeHandler()))
}

var lock = sync.RWMutex{}

// Unit Function For Rest
func CheckParameters(r *rest.Request, needed_fields []string) (bool, map[string]string) {
	var result map[string]string
	result = make(map[string]string)
	for _, field := range needed_fields {
		values, _ := r.URL.Query()[field]
		if len(values) < 1 {
			return false, result
		}
		result[field] = values[0]
	}
	return true, result
}

// DB Connection
var RConNum = 50
var sessions = createConnections(RConNum, "i7")
var sessions_taipei = createConnections(RConNum, "taipei_server")
var redis_client, r_status = redis_lib.NewClient()

// REST APIs
func Ping(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	w.WriteJson(bson.M{"Code": 0, "Result": "pong"})
	lock.RUnlock()
}

func PostUserEvent(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()

	// Check Header and mode
	db_name := "infohub_sandbox"
	if r.Header.Get("mode") == "production" {
		db_name = "infohub"
	}

	// Variable
	msg := "Error Message"
	Code := -1

	// Dealing with Post Body
	user_event := user_event.UserEvent{}
	err := r.DecodeJsonPayload(&user_event)
	if err != nil {
		log.Print(err)
		msg = err.Error()
	} else {
		// Session
		random_index := rand.Intn(20)
		session := sessions[random_index]

		// User Event
		_, msg = user_event.InsertOne(db_name, session)

		// User
		user := infohub_user.InfohubUser{Gaid: user_event.Gaid}
		user.Update(db_name, session, user_event.Content_id)
		Code = 0
	}

	// Return
	w.WriteJson(bson.M{"Code": Code, "Result": bson.M{"Message": msg}})
	lock.RUnlock()
}

func GetNewsByKeyword(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	needed_fields := []string{"keyword"}
	_, params := CheckParameters(r, needed_fields)
	n := news.News{}

	// Get news
	var results []news.News
	random_index := rand.Intn(RConNum)
	session := sessions[random_index]
	n.GetByKeyword(params["keyword"], session, &results)

	// Respond
	result := Result{"OK", nil, nil, nil}

	if len(results) > 0 {
		result = Result{"OK", results, nil, nil}
	}
	var respond = Respond{0, result}
	w.WriteJson(&respond)
	lock.RUnlock()
}

func GetStockList(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	n := stock.Stock{}

	// Get news
	var results []stock.Stock
	random_index := rand.Intn(RConNum)
	session := sessions[random_index]
	n.GetStockList(session, &results)

	w.WriteJson(bson.M{"Code": 0, "Result": bson.M{"Stocks": results}})
	lock.RUnlock()
}

func GetAll(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	needed_fields := []string{"country", "language", "category", "news_limit", "video_limit", "image_limit"}
	status, params := CheckParameters(r, needed_fields)
	random_index := rand.Intn(RConNum)
	if !status {
		var r_json map[string]string
		r_json = make(map[string]string)
		r_json["Status"] = "Error"
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteJson(&r_json)
	} else {
		image_limit, _ := strconv.Atoi(params["image_limit"])
		image_results := mongo_lib.GetImages(params["country"], params["language"], params["category"], sessions_taipei[random_index], image_limit, redis_client, r_status)
		video_limit, _ := strconv.Atoi(params["video_limit"])
		video_results := mongo_lib.GetVideos(params["country"], params["language"], params["category"], sessions_taipei[random_index], video_limit, redis_client, r_status)
		news_limit, _ := strconv.Atoi(params["news_limit"])
		news_results := mongo_lib.GetNews(params["country"], params["language"], params["category"], sessions[random_index], news_limit, redis_client, r_status)

		var result = Result{"No Data", nil, nil, nil}
		result = Result{"OK", news_results, video_results, image_results}

		var respond = Respond{0, result}

		w.WriteJson(&respond)
	}
	lock.RUnlock()
}

func GetImage(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	needed_fields := []string{"country", "language", "category"}
	status, params := CheckParameters(r, needed_fields)

	random_index := rand.Intn(RConNum)

	if !status {
		var r_json map[string]string
		r_json = make(map[string]string)
		r_json["Status"] = "Error"
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteJson(&r_json)
	} else {
		results := mongo_lib.GetImages(params["country"], params["language"], params["category"], sessions_taipei[random_index], 10, redis_client, r_status)
		var result = Result{"No Images", nil, nil, nil}

		if len(results) > 0 {
			result = Result{"OK", nil, nil, results}
		}

		var respond = Respond{0, result}

		w.WriteJson(&respond)
	}
	lock.RUnlock()
}

func GetVideo(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	needed_fields := []string{"country", "language", "category"}
	status, params := CheckParameters(r, needed_fields)

	random_index := rand.Intn(RConNum)

	if !status {
		var r_json map[string]string
		r_json = make(map[string]string)
		r_json["Status"] = "Error"
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteJson(&r_json)
	} else {
		results := mongo_lib.GetVideos(params["country"], params["language"], params["category"], sessions_taipei[random_index], 10, redis_client, r_status)
		var result = Result{"No Videos", nil, nil, nil}

		if len(results) > 0 {
			result = Result{"OK", nil, results, nil}
		}

		var respond = Respond{0, result}

		w.WriteJson(&respond)
	}
	lock.RUnlock()
}

func GetNews(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	needed_fields := []string{"country", "language", "category"}
	status, params := CheckParameters(r, needed_fields)

	random_index := rand.Intn(RConNum)

	if !status {
		var r_json map[string]string
		r_json = make(map[string]string)
		r_json["Status"] = "Error"
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteJson(&r_json)
	} else {
		results := mongo_lib.GetNews(params["country"], params["language"], params["category"], sessions[random_index], 10, redis_client, r_status)
		var result = Result{"No News", nil, nil, nil}

		if len(results) > 0 {
			result = Result{"OK", results, nil, nil}
		}

		var respond = Respond{0, result}

		w.WriteJson(&respond)
	}
	lock.RUnlock()
}

func createConnections(num int, account string) [50]*mgo.Session {
	var sessions [50]*mgo.Session
	for i := 0; i < num; i++ {
		mongo_url, err := MongoAccount(account)
		session, err := mgo.Dial(mongo_url)

		if err != nil {
			fmt.Printf("Error for Connection: %s\n", err)
		} else {
			fmt.Print("Connection Status: Good!\n")
		}
		sessions[i] = session
	}

	return sessions
}

func MongoAccount(_type string) (string, error) {
	var account Account
	filename, _ := filepath.Abs("mongo_lib/key.yaml")
	yamlFile, err := ioutil.ReadFile(filename)
	err = yaml.Unmarshal(yamlFile, &account)

	if err != nil {
		fmt.Println(err)
	}

	for _, v := range account.Mongo_users {
		if v.Type == _type {
			return fmt.Sprintf("mongodb://%s:%s@%s:%d", v.User, v.Password, v.Host, v.Port), nil
		}
	}
	return "", errors.New("User not found!")
}
