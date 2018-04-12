package mongo_lib

import (
	"crypto/sha1"
	"math/rand"
	"strings"

	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/go-redis/redis"

	"crypto/md5"

	"github.com/Terryhung/infohub_rest/gifimage"
	"github.com/Terryhung/infohub_rest/news"
	"github.com/Terryhung/infohub_rest/redis_lib"
	"github.com/Terryhung/infohub_rest/utils"
	"github.com/Terryhung/infohub_rest/video"
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

func RandomChoiceGeneral(dataset interface{}, _size int) {
	switch t := dataset.(type) {
	case *[]news.News:
		results := []news.News{}
		for i := 0; i < _size; i++ {
			s := *t
			random_index := rand.Intn(len(s))
			news := s[random_index]
			news.Append()
			results = append(results, news)
		}
		*t = results
	}
}

func RandomChoice(dataset []news.News, _size int) []news.News {
	results := []news.News{}
	for i := 0; i < _size; i++ {
		random_index := rand.Intn(len(dataset))
		news := dataset[random_index]
		news.Append()
		results = append(results, news)
	}
	return results
}

func RandomChoiceVideo(dataset []video.Video, _size int) []video.Video {
	results := []video.Video{}
	h := sha1.New()
	hasher := md5.New()
	for i := 0; i < _size; i++ {
		random_index := rand.Intn(len(dataset))
		dataset[random_index].ClassName = "video"
		dataset[random_index].Id = utils.MD5SHA1(dataset[random_index].Link, h, hasher)[:24]
		results = append(results, dataset[random_index])
	}
	return results
}

func RandomChoiceImage(dataset []gifimage.GifImage, _size int) []gifimage.GifImage {
	results := []gifimage.GifImage{}
	h := sha1.New()
	hasher := md5.New()
	for i := 0; i < _size; i++ {
		random_index := rand.Intn(len(dataset))
		dataset[random_index].ClassName = "image"
		dataset[random_index].Id = utils.MD5SHA1(dataset[random_index].Image_url, h, hasher)[:24]
		results = append(results, dataset[random_index])
	}
	return results
}

func InsertData(db_name string, col_name string, session *mgo.Session, data interface{}) (bool, string) {
	status := true
	msg := "Insert Data Successfully!"
	col := session.DB(db_name).C(col_name)
	err := col.Insert(data)
	if err != nil {
		status = false
		log.Fatal(err)
		msg = err.Error()
	}
	return status, msg
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

func GetImages(country string, language string, category string, session *mgo.Session, _size int, r_client *redis.Client, r_status bool) []gifimage.GifImage {
	var results []gifimage.GifImage
	if strings.Contains(category, "for") && strings.Contains(category, "you") {
		category = "headline"
	}

	keys := []string{country, language, category, "image"}
	key := strings.Join(keys, "-")
	if r_status {
		redis_lib.CheckExists(r_client, key, &results)
		if len(results) > 0 {
			results = RandomChoiceImage(results, _size)
			return results
		}
	}

	col := session.DB("droi").C("cache")
	constr := bson.M{"upserted_datetime": bson.M{"$gte": utils.NowTSNorm()*1000 - 86400000}, "category": category, "language": language, "country_array": country, "_from": bson.M{"$regex": "images/.*"}}
	_ = col.Find(constr).Limit(200).Sort("-upserted_datetime").All(&results)
	if len(results) == 0 {
		constr := bson.M{"upserted_datetime": bson.M{"$gte": utils.NowTSNorm()*1000 - 86400000}, "category": category, "language": language, "country_array": "ALL", "_from": bson.M{"$regex": "images/.*"}}
		_ = col.Find(constr).Limit(200).Sort("-upserted_datetime").All(&results)
	}

	if len(results) == 0 && language != "ar" && language != "in" {
		constr := bson.M{"upserted_datetime": bson.M{"$gte": utils.NowTSNorm()*1000 - 86400000}, "category": category, "language": "en", "country_array": "ALL", "_from": bson.M{"$regex": "images/.*"}}
		_ = col.Find(constr).Limit(200).Sort("-upserted_datetime").All(&results)
	}

	redis_lib.SetValue(r_client, key, results, 6000)

	if len(results) > 0 {
		results = RandomChoiceImage(results, _size)
	}
	return results
}

func GetVideos(country string, language string, category string, session *mgo.Session, _size int, r_client *redis.Client, r_status bool) []video.Video {
	var results []video.Video
	if strings.Contains(category, "for") && strings.Contains(category, "you") {
		category = "headline"
	}

	keys := []string{country, language, category, "video"}
	key := strings.Join(keys, "-")
	if r_status {
		redis_lib.CheckExists(r_client, key, &results)
		if len(results) > 0 {
			results = RandomChoiceVideo(results, _size)
			return results
		}
	}

	col := session.DB("droi").C("cache")
	constr := bson.M{"upserted_datetime": bson.M{"$gte": utils.NowTSNorm()*1000 - 86400000}, "category": category, "language": language, "country_array": country, "_from": bson.M{"$regex": "videos/.*"}}
	_ = col.Find(constr).Limit(200).Sort("-upserted_datetime").All(&results)
	if len(results) == 0 {
		constr := bson.M{"upserted_datetime": bson.M{"$gte": utils.NowTSNorm()*1000 - 86400000}, "category": category, "language": language, "country_array": "ALL", "_from": bson.M{"$regex": "videos/.*"}}
		_ = col.Find(constr).Limit(200).Sort("-upserted_datetime").All(&results)
	}

	if len(results) == 0 && language != "ar" && language != "in" {
		constr := bson.M{"upserted_datetime": bson.M{"$gte": utils.NowTSNorm()*1000 - 86400000}, "category": category, "language": "en", "country_array": "ALL", "_from": bson.M{"$regex": "videos/.*"}}
		_ = col.Find(constr).Limit(200).Sort("-upserted_datetime").All(&results)
	}

	redis_lib.SetValue(r_client, key, results, 6000)

	if len(results) > 0 {
		results = RandomChoiceVideo(results, _size)
	}
	return results
}

func GetNews(country string, language string, category string, session *mgo.Session, _size int, r_client *redis.Client, r_status bool) []news.News {
	var results []news.News
	keys := []string{country, language, category}
	key := strings.Join(keys, "-")

	if r_status {
		redis_lib.CheckExists(r_client, key, &results)
		if len(results) > 0 {
			results = RandomChoice(results, _size)
			return results
		}
	}

	if strings.Contains(category, "for") && strings.Contains(category, "you") {
		return GetForyou(country, language, category, session, _size, r_client, r_status)
	}

	col := session.DB("analysis").C("news_meta_baas")
	constr := bson.M{"source_date_int": bson.M{"$gte": utils.NowTSNorm() - 86400}, "category": category, "language": language, "country": country}
	_ = col.Find(constr).Limit(200).Sort("-source_date_int").All(&results)
	if len(results) == 0 {
		constr := bson.M{"source_date_int": bson.M{"$gte": utils.NowTSNorm() - 86400*3}, "category": category, "language": language, "country_array": bson.M{"$in": []string{"ALL", country}}}
		_ = col.Find(constr).Limit(200).Sort("-source_date_int").All(&results)
	}

	if len(results) == 0 && language != "ar" && language != "in" {
		constr := bson.M{"source_date_int": bson.M{"$gte": utils.NowTSNorm() - 86400*3}, "category": category, "language": "en", "country_array": bson.M{"$in": []string{"ALL", country}}}
		_ = col.Find(constr).Limit(200).Sort("-source_date_int").All(&results)
	}

	redis_lib.SetValue(r_client, key, results, 6000)

	if len(results) > 0 {
		results = RandomChoice(results, _size)
		// RandomChoiceGeneral(&results, _size)
	}

	return results
}

func CheckExist(db_name string, col_name string, session *mgo.Session, cond bson.M) bool {
	exists := false
	col := session.DB(db_name).C(col_name)
	count, _ := col.Find(cond).Count()
	if count > 0 {
		exists = true
	}
	return exists
}

func FindOne(db_name string, col_name string, session *mgo.Session, cond bson.M, data interface{}) bool {
	exists := false
	col := session.DB(db_name).C(col_name)
	count, _ := col.Find(cond).Count()
	if count > 0 {
		col.Find(cond).One(&data)
	}
	return exists
}
