package models

type Item struct {
	Id     int8   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"` // artistを含む
	Code   string `json:"code"`
	Place  string `json:"place"`
	Note   string `json:"note"`
	Image  []byte `json:"image"`
}
