package user_event

type UserEvent struct {
	Event_name string `json:"event_name"`
	Info_id    string `json:"info_id"`
	Ga_id      string `json:"ga_id"`
	News_id    string `json:"news_id"`
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
	if c.Event_name == "" || c.Ga_id == "" {
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
		}
	}
	return status
}
