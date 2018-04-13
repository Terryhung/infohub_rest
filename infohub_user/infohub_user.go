package infohub_user

import (
	"fmt"

	"github.com/Terryhung/infohub_rest/mongo_lib"
	"github.com/Terryhung/infohub_rest/news"
	"github.com/Terryhung/infohub_rest/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type InfohubUser struct {
	Device_id    string         `json:"device_id"`
	Gaid         string         `json:"gaid"`
	Country      string         `json:"country"`
	Favorite_cat map[string]int `json:"favorite_cat"`
	Last_login   int            `json:"last_login"`
}

func (i *InfohubUser) NewOne(db_name string, session *mgo.Session) (bool, string) {
	i.Last_login = utils.NowTS()
	status, msg := mongo_lib.InsertData(db_name, "user_profile", session, &i)
	return status, msg
}

func (i *InfohubUser) UpdateAttr(db_name string, session *mgo.Session, category []string) bool {
	status := false
	// condition [last login, favorite category]:
	// last login:
	last_login := utils.NowTS()

	// favorite category:
	favorite_cat := i.Favorite_cat
	for _, cat := range category {
		_, ok := favorite_cat[cat]
		if ok {
			favorite_cat[cat] += 1
		} else {
			if len(favorite_cat) == 0 {
				favorite_cat = make(map[string]int)
			}
			favorite_cat[cat] = 1
		}
	}

	// Merge condition
	cond := bson.M{"last_login": last_login, "favorite_cat": favorite_cat}
	_id := bson.M{"gaid": i.Gaid}

	// Call Mongo API
	status = mongo_lib.UpdateOne(db_name, "user_profile", session, cond, _id)
	return status
}

func (i *InfohubUser) Update(db_name string, session *mgo.Session, news_id string) {
	// Check user exists or not.
	var user_profile = bson.M{"gaid": i.Gaid}
	exists, col := mongo_lib.FindOne(db_name, "user_profile", session, user_profile, i)

	// If user does not exist, create new one. Else find one.
	if !exists {
		fmt.Print("User Not Exists!\n")
		i.NewOne(db_name, session)
	} else {
		col.Find(user_profile).One(&i)
	}

	// Start to Update User Data by News.
	news := news.News{Id: news_id}
	news_exists := news.CheckExist(session)
	if news_exists {
		i.UpdateAttr(db_name, session, news.Category)
	}
}
