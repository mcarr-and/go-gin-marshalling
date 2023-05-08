package model

type Album struct {
	ID     int     `json:"id" binding:"min=1,max=10000"`
	Title  string  `json:"title" binding:"required,min=2,max=1000"`
	Artist string  `json:"artist" binding:"required,min=2,max=1000"`
	Price  float64 `json:"price" binding:"required,min=0.0,max=10000.00"`
}
