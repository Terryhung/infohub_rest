package user_event

import (
	"github.com/Terryhung/infohub_rest/mongo_lib"
	"github.com/Terryhung/infohub_rest/utils"
	mgo "gopkg.in/mgo.v2"
)

type UserEvent struct {
	Event             string `json:"event"`
	Event_src         string `json:"event_src"`
	Dev_id            string `json:"dev_id"`
	Gaid              string `json:"gaid"`
	Country           string `json:"country"`
	Area              string `json:"area"`
	Content_id        string `json:"content_id"`
	Content_type      string `json:"content_type"`
	Content_category  string `json:"content_category"`
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
		"view":     0,
		"add":      0,
		"like":     1,
		"dislike":  1,
		"share":    1,
		"comment":  1,
		"bookmark": 1,
	}
	// Check Needed fields: Can not be nil
	if c.Event == "" || c.Gaid == "" {
		status = false
	}

	// Check Event Name Valid or not
	check_type, ok := valid_event_name[c.Event]
	if !ok {
		status = false
	} else {
		switch check_type {
		// News id cant not be nil
		case 1:
			if c.Content_id == "" {
				status = false
			}
			// Must provide keyword string
		}
	}
	return status
}
