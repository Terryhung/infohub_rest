package mongo_lib

import (
	"crypto/sha1"
	"math"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/go-redis/redis"

	"crypto/md5"
	"encoding/hex"

	"github.com/Terryhung/infohub_rest/news"
	"github.com/Terryhung/infohub_rest/redis_lib"
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

func RandomChoice(dataset []news.News, _size int) []news.News {
	results := []news.News{}
	h := sha1.New()
	hasher := md5.New()
	for i := 0; i < _size; i++ {
		random_index := rand.Intn(len(dataset))
		dataset[random_index].ClassName = "news"
		hasher.Write([]byte(dataset[random_index].Link))
		_id := hex.EncodeToString(hasher.Sum(nil))
		h.Write([]byte(_id))
		bs := hex.EncodeToString(h.Sum(nil))
		dataset[random_index].Id = bs[:24]
		results = append(results, dataset[random_index])
	}
	return results
}

func GetForyou(country string, language string, category string, session *mgo.Session, _size int, r_client *redis.Client, r_status bool) []news.News {
	categories := []string{"pets", "girls", "food"}
	results := []news.News{}
	h_size := _size / 2

	for len(results) < h_size && len(categories) > 0 {
		random_index := rand.Intn(len(categories))
		category := categories[random_index]
		p_results := GetNews(country, language, category, session, 1, r_client, r_status)
		if len(p_results) > 0 {
			results = append(results, p_results[0])
		} else {
			categories = append(categories[:random_index], categories[random_index+1:]...)
		}
	}

	headline_result := GetNews(country, language, "headline", session, _size-h_size, r_client, r_status)
	for i := 0; i < len(headline_result); i++ {
		results = append(results, headline_result[i])
	}

	return results
}

func GetNews(country string, language string, category string, session *mgo.Session, _size int, r_client *redis.Client, r_status bool) []news.News {
	var results []news.News
	keys := []string{country, language, category}
	key := strings.Join(keys, "-")

	if r_status {
		_results := redis_lib.CheckExists(r_client, key)
		if len(_results) > 0 {
			_results = RandomChoice(_results, _size)
			return _results
		}
	}

	if strings.Contains(category, "for") && strings.Contains(category, "you") {
		return GetForyou(country, language, category, session, _size, r_client, r_status)
	}

	col := session.DB("analysis").C("news_meta_baas")
	constr := bson.M{"source_date_int": bson.M{"$gte": NowTSNorm() - 86400}, "category": category, "language": language, "country": country}
	_ = col.Find(constr).Limit(200).Sort("-source_date_int").All(&results)
	if len(results) == 0 {
		constr := bson.M{"source_date_int": bson.M{"$gte": NowTSNorm() - 86400*3}, "category": category, "language": language, "country_array": "ALL"}
		_ = col.Find(constr).Limit(200).Sort("-source_date_int").All(&results)
	}

	redis_lib.SetValue(r_client, key, results, 6000)

	if len(results) > 0 {
		results = RandomChoice(results, _size)
	}

	return results
}
