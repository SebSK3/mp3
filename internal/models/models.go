package internal

type Song struct {
	Id       int64
	Url      string `form:"Url"`
	Mix      bool   `form:"Mix"`
	Artist   string `form:"Artist,omitempty"`
	Title    string `form:"Title,omitempty"`
	Album    string `form:"Album,omitempty"`
	Fullname string
	Filename string
	Progress int64
}
