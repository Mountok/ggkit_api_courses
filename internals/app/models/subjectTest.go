package models

type SubjectTest struct {
	Id        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	SubjectId int    `json:"subject_id" db:"subject_id"`
}

type TestQuestion struct {
	Id       int    `json:"id" db:"id"`
	Question string `json:"question" db:"question"`
	Options  string `json:"options" db:"options"`
	TestId   int    `json:"test_id" db:"test_id"`
}

type CompletedTest struct {
	Id     int    `json:"id" db:"id"`
	TestId int    `json:"test_id" db:"test_id"`
	UserId string `json:"user_id" db:"user_id"`
	Points int    `json:"points" db:"points"`
}
type CompletedTestCheck struct {
	Id        int    `json:"id" db:"id"`
	TestId    int    `json:"test_id" db:"test_id"`
	SubjectId int    `json:"subject_id" db:"subject_id"`
	UserId    string `json:"user_id" db:"user_id"`
	Points    int    `json:"points" db:"points"`
}

type QuestionCheckReq struct {
	QuestionId int    `json:"question_id"`
	Answer     string `json:"answer"`
}
