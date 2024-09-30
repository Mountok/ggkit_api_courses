package models

type User struct {
	Id         string `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	CreateDate string `json:"create_date"`
}

type ValidateUser struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}

type RatingUser struct {
	Id       string `json:"id" db:"user_id"`
	Username string `json:"username" db:"full_name"`
	ImageUrl string `json:"image" db:"image"`
	Score    int    `json:"score" db:"score"`
}

type LastSubject struct {
	Id         int    `json:"id" db:"id"`
	UserId     string `json:"user_id" db:"user_id"`
	SubjectsId int    `json:"subjects_id" db:"subjects_id"`
}
