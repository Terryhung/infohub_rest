package video

type Video struct {
	Title           string   `json:"title"`
	Source_name     string   `json:"source_name"`
	Image_url       string   `json:"image_url"`
	Like_numbers    int      `json:"like_numbers"`
	Unlike_numbers  int      `json:"unlike_numbers"`
	Description     string   `json:"description"`
	Link            string   `json:"link"`
	Source_date_int int      `json:"source_date_int"`
	Similar_ids     []string `json:"similar_ids"`
	Video_length    int      `json:"video_length"`
	ClassName       string   `json:"_ClassName"`
	Id              string   `json:"_Id"`
}
