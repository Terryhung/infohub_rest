package mongo_lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Terryhung/infohub_rest/news"
)

const TBLURL = "https://contentapi.celltick.com/mediaApi/v1.0/personal/content?publisherId=JC_InfohubLegacy-Web&key=x8fPbq6FRUPD5DUOYxOTBkipjjuztcB4"

type TBLResp struct {
	Content []TBL `json:"content"`
}

type TBL struct {
	Link         string      `json:"contentURL"`
	Title        string      `json:"title"`
	Updated      int         `json:"publishedAt"`
	ChannelName  string      `json:"contentSourceDisplay"`
	ChannelImage string      `json:"contentSourceLogo,omitempty"`
	ImagesObj    TBLImageObj `json:"images"`
	Images       []string    `json:"clearImages"`
	Description  string      `json:"summary"`
}

type TBLImageObj struct {
	MainImage   TBLImage `json:"mainImage"`
	MainImageTB TBLImage `json:"mainImageThumbnail"`
}

type TBLImage struct {
	URL string `json:"url"`
}

func (m *TBL) toNews() (news news.News) {
	news.Title = m.Title
	news.Link = m.Link
	news.Description = m.Description
	news.Image_url_array = m.Images
	news.Source_name = "Taboola"

	return
}

func QueryTBLNews(cty, lang, gaid string, limit int) (news []news.News) {
	url := fmt.Sprintf("%s&countryCode=%s&language=%s&limit=%d&userId=%s", TBLURL, cty, lang, limit, gaid)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("[QueryMSAD] fail to get Taboola AD, err: %+v", err)
		return
	}

	defer resp.Body.Close()

	TBLObj := TBLResp{}

	json.NewDecoder(resp.Body).Decode(&TBLObj)

	ads := TBLObj.Content
	for _, ad := range ads {
		images := []string{}
		images = append(images, ad.ImagesObj.MainImage.URL)
		images = append(images, ad.ImagesObj.MainImageTB.URL)
		ad.Images = images
		news = append(news, ad.toNews())
	}

	return
}
