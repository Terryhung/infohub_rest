package mongo_lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Terryhung/infohub_rest/news"
)

const (
	TBLURL      = "https://api.taboola.com/1.2/json"
	INFOHUB_URL = "http://www.infohubapp.com/%s-%s/index.html"
)

var TBLCTY = map[string][]string{}

type ClientInfo struct {
	PublisherId string
	Key         string
	Path        string
}

var TBChannelMapping = map[string]ClientInfo{
	"DE": ClientInfo{PublisherId: "1259378", Key: "76fa80afb7d21bb2b30dc665d1cb707812fae598", Path: "ume-infohub-sc-germany"},
	"IT": ClientInfo{PublisherId: "1259377", Key: "76fa80afb7d21bb2b30dc665d1cb707812fae598", Path: "ume-infohub-sc-italy"},
	"RU": ClientInfo{PublisherId: "1259375", Key: "76fa80afb7d21bb2b30dc665d1cb707812fae598", Path: "ume-infohub-sc-russia"},
	"ES": ClientInfo{PublisherId: "1259376", Key: "76fa80afb7d21bb2b30dc665d1cb707812fae598", Path: "ume-infohub-sc-spain"},
}

type TBLResp struct {
	Content []TBL `json:"list"`
}

type TBIMG struct {
	Link string `json:"url"`
}

type TBL struct {
	Link        string  `json:"url"`
	Title       string  `json:"name"`
	ChannelName string  `json:"branding"`
	ImageURL    []TBIMG `json:"thumbnail"`
	Description string  `json:"promotedText"`
}

func (m *TBL) toNews() (news news.News) {
	news.Title = m.Title
	news.Link = m.Link
	news.Description = m.Description
	imgs := []string{}
	for _, img := range m.ImageURL {
		imgs = append(imgs, img.Link)
	}
	news.Image_url_array = imgs
	news.Source_name = "Taboola"

	return
}

func QueryTBLNews(cty, lang, gaid string, limit int) (news []news.News) {
	mapping, ok := TBChannelMapping[cty]
	if !ok {
		return
	}
	infohub := fmt.Sprintf(INFOHUB_URL, lang, cty)
	url := fmt.Sprintf("%s/%s/recommendations.get?app.type=mobile&app.apikey=%s&placement.rec-count=%d&placement.organic-type=text&user.session=init&source.type=text&source.id=%s&source.url=%s", TBLURL, mapping.Path, mapping.Key, limit, mapping.PublisherId, infohub)
	resp, err := http.Get(url)
	log.Printf(url)

	if err != nil {
		log.Printf("[QueryMSAD] fail to get Taboola AD, err: %+v", err)
		return
	}

	defer resp.Body.Close()

	TBLObj := TBLResp{}

	json.NewDecoder(resp.Body).Decode(&TBLObj)

	ads := TBLObj.Content
	for _, ad := range ads {
		news = append(news, ad.toNews())
	}

	return
}
