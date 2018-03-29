package gifimage

type GifImage struct {
	Title           string `json:"title"`
	Source_name     string `json:"source_name"`
	Image_url       string `json:"image_url"`
	Like_numbers    int    `json:"like_numbers"`
	Unlike_numbers  int    `json:"unlike_numbers"`
	Description     string `json:"description"`
	Source_date_int int    `json:"source_date_int"`
	ClassName       string `json:"_ClassName"`
	Id              string `json:"_Id"`
	Video_url       string `json:"video_url"`
}
