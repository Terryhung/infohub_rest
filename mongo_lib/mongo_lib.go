package mongo_lib

import (
	"math"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	Title             string   `json:"title"`
	Source_name       string   `json:"source_name"`
	Image_url_array   []string `json:"image_url_array"`
	Image_url         string   `json:"image_url"`
	Like_numbers      int      `json:"link_numbers"`
	Unlike_numbers    int      `json:"unlink_numbers"`
	Description       string   `json:"description"`
	Page_link         string   `json:"page_link"`
	Explicit_keywords []string `json:"explicit_keywords"`
	Source_date       string   `json:"source_date"`
	Similar_ids       []string `json:"similar_ids"`
	ClassName         string   `json:"_ClassName"`
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

func RandomChoice(dataset []News, _size int) []News {
	results := []News{}
	for i := 0; i < _size; i++ {
		random_index := rand.Intn(len(dataset))
		dataset[random_index].ClassName = "news"
		results = append(results, dataset[random_index])
	}
	return results
}

func GetForyou(country string, language string, category string, session *mgo.Session, _size int) []News {
	categories := []string{"pets", "girls", "food"}
	results := []News{}
	h_size := _size / 2

	for len(results) < h_size {
		random_index := rand.Intn(len(categories))
		category := categories[random_index]
		p_results := GetNews(country, language, category, session, 1)
		if len(p_results) > 0 {
			results = append(results, p_results[0])
		}
	}

	headline_result := GetNews(country, language, "headline", session, _size-h_size)
	for i := 0; i < len(headline_result); i++ {
		results = append(results, headline_result[i])
	}

	return results
}

func GetNews(country string, language string, category string, session *mgo.Session, _size int) []News {
	var results []News

	if strings.Contains(category, "for") && strings.Contains(category, "you") {
		return GetForyou(country, language, category, session, _size)
	}

	col := session.DB("analysis").C("news_meta_baas")
	constr := bson.M{"source_date_int": bson.M{"$gte": NowTSNorm() - 86400}, "category": category, "language": language, "country": country}
	_ = col.Find(constr).Limit(200).Sort("-source_date_int").All(&results)
	if len(results) == 0 {
		constr := bson.M{"source_date_int": bson.M{"$gte": NowTSNorm() - 86400*3}, "category": category, "language": language, "country_array": bson.M{"$in": []string{"ALL", country}}}
		_ = col.Find(constr).Limit(200).Sort("-source_date_int").All(&results)
	}

	if len(results) > 0 {
		results = RandomChoice(results, _size)
	}

	return results
}
