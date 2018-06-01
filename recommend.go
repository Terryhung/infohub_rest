package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/Terryhung/infohub_rest/mongo_lib"
	"github.com/Terryhung/infohub_rest/news"
	"github.com/Terryhung/infohub_rest/redis_lib"
	"github.com/Terryhung/infohub_rest/utils"
	"github.com/go-redis/redis"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type InfohubUser struct {
	Gaid   string `json:"gaid" bson:"gaid"`
	Cand   []int  `json:"candidates" bson:"candidates"`
	Method string `json:"method" bson:"method"`
}

func Recommendar(gaid string, lang string, cty string, session *mgo.Session, r_client *redis.Client, r_status bool) (bool, []news.News) {
	db := "analysis"
	col := "implicit_recommendation_news"
	user := InfohubUser{}
	news_results := []news.News{}
	status, _ := mongo_lib.FindOne(db, col, session, bson.M{"gaid": gaid}, &user)
	if !status {
		// Check general user
		g_gaids := []string{"general-user", cty}
		g_gaid := strings.Join(g_gaids, "-")

		// Try
		status, _ = mongo_lib.FindOne(db, col, session, bson.M{"gaid": g_gaid}, &user)
		if !status {
			return false, news_results
		}
	}

	// Query News
	news_db := "analysis"
	news_col := "news_meta_baas"

	for _, c := range user.Cand[:10] {
		cond := bson.M{"hier_category": c, "language": lang}
		var c_news []news.News

		// Concate Key for cache
		keys := []string{lang, strconv.Itoa(c)}
		key := strings.Join(keys, "-")
		need_cache := false

		// Check Cache
		if r_status {

			// Check data in redis
			redis_lib.CheckExists(r_client, key, &c_news)

			// No Data: check data in Mongo
			if len(c_news) == 0 {
				mongo_lib.Find(news_db, news_col, session, cond, &c_news)
				need_cache = true
			}
		} else {
			mongo_lib.Find(news_db, news_col, session, cond, &c_news)
		}

		// Random Pick
		if len(c_news) > 0 {
			random_index := rand.Intn(len(c_news))
			n := c_news[random_index]
			n.By = user.Method
			n.Id = utils.SpecialID(n.Link)
			news_results = append(news_results, n)
		}

		// Need Cache?
		if need_cache && r_status {
			redis_lib.SetValue(r_client, key, c_news, 3600)
		}
	}

	return status, news_results
}
