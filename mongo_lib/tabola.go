package mongo_lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Terryhung/infohub_rest/news"
)

const TBLURL = "https://contentapi.celltick.com/mediaApi/v1.0/mid/promoted/?publisherId=JC_InfohubLegacy-Web&key=x8fPbq6FRUPD5DUOYxOTBkipjjuztcB4"

type TBLResp struct {
	Content []TBL `json:"content"`
}

type TBL struct {
	Link        string `json:"actionUri"`
	Title       string `json:"title"`
	ChannelName string `json:"contentSource"`
	ImageURL    string `json:"imageUrl"`
	Description string `json:"promotedText"`
}

func (m *TBL) toNews() (news news.News) {
	news.Title = m.Title
	news.Link = m.Link
	news.Description = m.Description
	news.Image_url_array = []string{m.ImageURL}
	news.Source_name = "Taboola"

	return
}

func QueryTBLNews(cty, lang, gaid string, limit int) (news []news.News) {
	url := fmt.Sprintf("%s&countryCode=%s&language=%s&limit=%d&userId=%s", TBLURL, cty, lang, limit, gaid)
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
