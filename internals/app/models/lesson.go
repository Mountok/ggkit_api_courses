package models

type Lesson struct {
	Id      int    `json:"id"`
	Upkeep  string `json:"upkeep"`
	ThemeId int    `json:"theme_id"`
}

type DoneLesson struct {
	Id      int `json:"id"`
	ThemeId int `json:"theme_id"`
	UserId  int `json:"user_id"`
}
