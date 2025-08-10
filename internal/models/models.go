package internal

type Song struct {
	Id                  int64
	Url                 string `form:"Url"`
	Mix                 bool   `form:"Mix"`
	Artist              string `form:"Artist,omitempty"`
	Title               string `form:"Title,omitempty"`
	Album               string `form:"Album,omitempty"`
	CustomImageUrl      string `form:"ImgUrl,omitempty"`
	CustomFFmpegCommand string `form:"FFmpegCmd,omitempty"`
	Fullname            string
	Filename            string
	Progress            int64
}
