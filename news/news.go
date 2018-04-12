package news

import (
	"github.com/Terryhung/infohub_rest/utils"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type News struct {
	Title             string   `json:"title"`
	Source_name       string   `json:"source_name"`
	Image_url_array   []string `json:"image_url_array"`
	Image_url         string   `json:"image_url"`
	Like_numbers      int      `json:"like_numbers"`
	Unlike_numbers    int      `json:"unlike_numbers"`
	Description       string   `json:"description"`
	Page_link         string   `json:"page_link"`
	Link              string   `json:"link"`
	Explicit_keywords []string `json:"explicit_keywords"`
	Source_date_int   int      `json:"source_date_int"`
	Similar_ids       []string `json:"similar_ids"`
	ClassName         string   `json:"_ClassName"`
	Id                string   `json:"_Id"`
	Category          []string `json:"category"`
}

func (n *News) Append() {
	n.ClassName = "news"
	n.Id = utils.SpecialID(n.Link)
}

func (n *News) CheckExist(session *mgo.Session) bool {
	exists := false
	cond := bson.M{"_baas_id": n.Id}
	col := session.DB("analysis").C("news_meta_baas")
	count, _ := col.Find(cond).Count()
	if count > 0 {
		col.Find(cond).One(&n)
		exists = true
	}
	return exists
}
