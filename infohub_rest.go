package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/Terryhung/infohub_rest/mongo_lib"
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

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultCommonStack...)
	router, err := rest.MakeRouter(
		rest.Get("/get_news", GetNews),
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

var sessions = createConnections(20)

func GetNews(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	needed_fields := []string{"country", "language", "category"}
	status, params := CheckParameters(r, needed_fields)

	random_index := rand.Intn(10)

	if !status {
		var r_json map[string]string
		r_json = make(map[string]string)
		r_json["Status"] = "Error"
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.WriteJson(&r_json)
	} else {
		fmt.Printf("Session %+v\n", sessions[random_index])
		results := mongo_lib.GetNews(params["country"], params["language"], params["category"], sessions[random_index])
		w.WriteJson(&results)
	}
	lock.RUnlock()
}

func createConnections(num int) [20]*mgo.Session {
	var sessions [20]*mgo.Session
	for i := 0; i < num; i++ {
		mongo_url, err := MongoAccount("normal")
		session, err := mgo.Dial(mongo_url)

		if err != nil {
			fmt.Printf("Error for Connection: %s\n", err)
		}
		fmt.Printf("Session %+v\n", session)
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
