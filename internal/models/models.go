package models

type Item struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Author   string `json:"author"`
	Code     string `json:"code"`
	Purchase string `json:"purchase"`
	Place    string `json:"place"`
	Note     string `json:"note"`
	Image    string `json:"image"`
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

// ユーザー関連
type UserInternal struct {
	Id           string
	UserName     string
	PasswordHash string
	UUID         string
}

type UserExternalResponse struct {
	UserName string `json:"username"`
	UUID     string `json:"uuid"`
}

type UserExternalLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
