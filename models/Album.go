package models

type Album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title" binding:"required"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}
