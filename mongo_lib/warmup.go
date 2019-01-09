package mongo_lib

import (
	"github.com/Terryhung/infohub_rest/gifimage"
	"github.com/Terryhung/infohub_rest/news"
	"github.com/Terryhung/infohub_rest/redis_lib"
	"github.com/Terryhung/infohub_rest/utils"
	"github.com/go-redis/redis"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var DB = "analysis"
var IMG_COL = "image_cache"

func NewsQuery(country string, language string, category string, session *mgo.Session) []news.News {
	// Query News by (country, language, category)
	var news_result []news.News

	// Generate Query Condition
	offset := 1
	col_name := "news_meta_baas"

	// Check is special category
	_, ok := Special_Category[category]
	if ok {
		offset = 6
		col_name = "news_meta"
	}

	date_cond := utils.NowTSNorm() - 86400*int32(offset)
	cond := GenCondition(category, language, country, bson.M{}, date_cond)

	// Initial session
	col := session.DB(DB).C(col_name)

	// Sort by date
	col.Find(cond).Sort("-source_date_int").Limit(200).All(&news_result)

	return news_result
}

func ImageQuery(country string, language string, category string, session *mgo.Session) []gifimage.GifImage {
	// Query Image by (country, language, category)
	var img_result []gifimage.GifImage

	// Generate Query Condition
	offset := 1

	// Check is special category
	_, ok := Special_Category[category]
	if ok {
		offset = 6
	}

	date_cond := utils.NowTSNorm() - 86400*int32(offset)
	cond := bson.M{"upserted_datetime": bson.M{"$gte": date_cond}, "category": category, "language": language}

	// Initial session
	col := session.DB(DB).C(IMG_COL)

	// Sort by date
	col.Find(cond).Sort("-upserted_datetime").Limit(200).All(&img_result)

	return img_result
}

func Warmup(country string, language string, category string, session *mgo.Session, r_client *redis.Client) bool {
	// Warm Up Redis Cache
	status := false

	// Query result
	news_result := NewsQuery(country, language, category, session)
	img_result := ImageQuery(country, language, category, session)

	if len(news_result) >= 100 {
		// Initial Redis Key
		news_rds_key := GenRedisKey([]string{country, language, category})

		// Warmup cache
		redis_lib.SetValue(r_client, news_rds_key, news_result, 900)

		// Set status == true
		status = true
	}

	if len(img_result) >= 60 {
		// Initial Redis Key
		img_rds_key := GenRedisKey([]string{country, language, category, "image"})

		// Warmup cache
		redis_lib.SetValue(r_client, img_rds_key, img_result, 900)

		// Set status == true
		status = true
	}

	return status
}
