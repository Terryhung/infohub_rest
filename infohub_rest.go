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

	"github.com/Terryhung/infohub_rest/mongo_lib"
	"github.com/Terryhung/infohub_rest/redis_lib"
	"github.com/ant0ine/go-json-rest/rest"

	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"
)

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

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	api.Use(&rest.GzipMiddleware{})
	router, err := rest.MakeRouter(
		rest.Get("/get_news", GetNews),
		rest.Get("/get_video", GetVideo),
		rest.Get("/get_image", GetImage),
		rest.Get("/get_all", GetAll),
	)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8787", api.MakeHandler()))
}

var lock = sync.RWMutex{}

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

var sessions = createConnections(20, "i7")
var sessions_taipei = createConnections(20, "taipei_server")

var redis_client, r_status = redis_lib.NewClient()

func GetAll(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	needed_fields := []string{"country", "language", "category", "news_limit", "video_limit", "image_limit"}
	status, params := CheckParameters(r, needed_fields)
	random_index := rand.Intn(20)
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

	random_index := rand.Intn(20)

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

	random_index := rand.Intn(20)

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

	random_index := rand.Intn(20)

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

func createConnections(num int, account string) [20]*mgo.Session {
	var sessions [20]*mgo.Session
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
