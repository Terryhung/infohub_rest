package user_event

import (
	"github.com/Terryhung/infohub_rest/mongo_lib"
	"github.com/Terryhung/infohub_rest/utils"
	mgo "gopkg.in/mgo.v2"
)

type UserEvent struct {
	Event_name        string `json:"event_name"`
	Info_id           string `json:"info_id"`
	Device_id         string `json:"device_id"`
	Gaid              string `json:"gaid"`
	Country           string `json:"country"`
	Area              string `json:"area"`
	News_id           string `json:"news_id"`
	Keyword           string `json:"keyword"`
	Created_timestamp int    `json:"created_timestamp"`
}

func (c *UserEvent) AppendField() {
	c.Created_timestamp = utils.NowTS()
}

func (c *UserEvent) InsertOne(db_name string, session *mgo.Session) (bool, string) {
	status := false
	msg := "User Event format Error!"
	if c.Check() {
		c.AppendField()
		status, msg = mongo_lib.InsertData(db_name, "user_event", session, &c)
	}
	return status, msg
}

func (c UserEvent) Check() bool {
	status := true
	valid_event_name := map[string]int{
		"click_profile":     0,
		"read_article":      1,
		"read_notification": 1,
		"comment":           1,
		"click_like":        1,
		"browse":            1,
		"click_keyword":     2,
		"search":            2,
		"add_category":      3,
	}
	// Check Needed fields: Can not be nil
	if c.Event_name == "" || c.Gaid == "" {
		status = false
	}

	// Check Event Name Valid or not
	check_type, ok := valid_event_name[c.Event_name]
	if !ok {
		status = false
	} else {
		switch check_type {
		// News id cant not be nil
		case 1:
			if c.News_id == "" {
				status = false
			}
		// Must provide keyword string
		case 2:
			if c.Keyword == "" {
				status = false
			}
		}
	}
	return status
}
