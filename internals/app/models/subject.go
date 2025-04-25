package models

type Subject struct {
	Id             int    `json:"id"`
	Title          string `json:"title"`
	Image          string `json:"image"`
	Description    string `json:"description"`
	Iscertificated string `json:"iscertificated"`
}

type SubjectResponse struct {
	Id               int    `json:"id"`
	Title            string `json:"title"`
	Image            string `json:"image"`
	Description      string `json:"description"`
	Iscertificated   string `json:"iscertificated"`
	Total_themes     int    `json:"total_themes"`
	Completed_themes int    `json:"completed_themes"`
}
