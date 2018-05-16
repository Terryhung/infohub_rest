package main

import (
	"math/rand"

	"github.com/Terryhung/infohub_rest/mongo_lib"
	"github.com/Terryhung/infohub_rest/news"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type InfohubUser struct {
	Gaid string `json:"gaid" bson:"gaid"`
	Top  []int  `json:"top" bson:"top"`
}

func Recommendar(gaid string, lang string, session *mgo.Session) (bool, []news.News) {
	db := "infohub_sandbox"
	col := "user_profile"
	user := InfohubUser{}
	news_results := []news.News{}

	status, _ := mongo_lib.FindOne(db, col, session, bson.M{"gaid": gaid}, &user)
	if !status {
		return false, news_results
	}

	// Query News
	news_db := "analysis"
	news_col := "news_meta_baas"

	for _, c := range user.Top {
		cond := bson.M{"hier_category": c, "language": lang}
		var c_news []news.News
		mongo_lib.Find(news_db, news_col, session, cond, &c_news)

		// Random Pick
		random_index := rand.Intn(len(c_news))
		n := c_news[random_index]
		n.By = "Hier"
		news_results = append(news_results, n)
	}

	return status, news_results
}
