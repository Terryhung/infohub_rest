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
	Device_id   string `json:"device_id"`
	Gaid        string `json:"gaid"`
	Country     string `json:"country"`
	FavoriteCAT bson.M `json:"favorite_cat"`
	Last_login  int    `json:"last_login"`
}

func (i *InfohubUser) NewOne(db_name string, session *mgo.Session) (bool, string) {
	i.Last_login = utils.NowTS()
	status, msg := mongo_lib.InsertData(db_name, "user_profile", session, &i)
	return status, msg
}

func (i *InfohubUser) Update(db_name string, session *mgo.Session, news_id string) {
	var user_profile = bson.M{"gaid": i.Gaid}
	exists := mongo_lib.CheckExist(db_name, "user_profile", session, user_profile)
	if !exists {
		fmt.Print("Not Exists!\n")
		_, msg := i.NewOne(db_name, session)
		fmt.Print(msg)
	}
	news := news.News{Id: news_id}
	news_exists := news.CheckExist(session)
	if news_exists {
		fmt.Print(news)
	}
}
