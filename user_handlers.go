package main

import (
	"math/rand"

	"github.com/ant0ine/go-json-rest/rest"
)

func GetForYou(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()

	// Need gaid, language, country
	needed_fields := []string{"gaid", "language"}
	_, params := CheckParameters(r, needed_fields)

	// Get Recommend News
	session := sessions[rand.Intn(RConNum)]
	status, news_results := Recommendar(params["gaid"], params["language"], session)

	if !status {
		GetNews(w, r)
		lock.RUnlock()
	} else {
		var resp = Result{"No News", nil, nil, nil}
		if len(news_results) > 0 {
			resp = Result{"OK", news_results, nil, nil}
		}
		var respond = Respond{0, resp}
		w.WriteJson(&respond)
		lock.RUnlock()
	}
}
