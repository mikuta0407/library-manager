package models

type Item struct {
	Id     int8
	Title  string
	Author string // artistを含む
	Code   string
	Place  string
	Note   string
	Image  []byte
}
