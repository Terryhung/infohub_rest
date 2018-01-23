package mongo_lib

import (
	"errors"
	"fmt"
	"math"
	"time"

	"io/ioutil"
	"path/filepath"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

type News struct {
	Title string
}

func NowTSNorm() int32 {
	ts := time.Now().Unix()
	return int32(math.Floor(float64(ts)/86400) * 86400)
}

func NowMonth() string {
	date_time := time.Now().Format("2006-01")
	return date_time
}

func NowDate() string {
	date_time := time.Now().Format("2006-01-01")
	return date_time
}

func GetNews(country string, language string, category string, session *mgo.Session) []News {
	var results []News

	col := session.DB("droi").C("cache")
	constr := bson.M{"source_date_int": bson.M{"$gte": NowTSNorm() - 86400}, "category": category, "language": language, "country": country}
	_ = col.Find(constr).Limit(10).All(&results)

	return results
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
