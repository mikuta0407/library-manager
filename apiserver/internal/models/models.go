package models

type Item struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Code     string `json:"code"`
	Purchase string `json:"purchase"`
	Place    string `json:"place"`
	Note     string `json:"note"`
	Image    []byte `json:"image"`
}

type ItemArray struct {
	ItemList []Item `json:"items"`
}

// 返答メッセージJSON用構造体
type SuccessResponseMessage struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

type ErrMessage struct {
	ErrMessage string `json:"message"`
	ErrDetail  string `json:"detail"`
}
