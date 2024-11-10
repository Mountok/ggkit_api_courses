package models

type Subject struct {
	Id             int    `json:"id"`
	Title          string `json:"title"`
	Image          string `json:"image"`
	Description    string `json:"description"`
	Iscertificated string `json:"iscertificated"`
}
