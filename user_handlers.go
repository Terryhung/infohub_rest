package main

import (
	"math/rand"
	"strconv"

	"github.com/Terryhung/infohub_rest/mongo_lib"
	"github.com/ant0ine/go-json-rest/rest"
)

func GetForYou(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()

	// Need gaid, language, country
	needed_fields := []string{"country", "language", "category", "news_limit", "video_limit", "image_limit", "gaid"}
	_, params := CheckParameters(r, needed_fields)

	// Get Recommend News
	session := sessions[rand.Intn(RConNum)]
	status, news_results := Recommendar(params["gaid"], params["language"], session)

	// Get Images
	image_limit, _ := strconv.Atoi(params["image_limit"])
	image_results := mongo_lib.GetImages(params["country"], params["language"], params["category"], sessions_taipei[rand.Intn(RConNum)], image_limit, redis_client, r_status)

	if !status || len(news_results) < 4 {
		GetAll(w, r)
		lock.RUnlock()
	} else {
		var resp = Result{"No News", nil, nil, nil}

		if len(news_results) > 0 {
			resp = Result{"OK", news_results, nil, image_results}
		}
		var respond = Respond{0, resp}
		w.WriteJson(&respond)
		lock.RUnlock()
	}
}
