package models


type Book struct {
	ID int `json:"id,omitempty"`
	Author string `json:"author"`
	Title string `json:"title"`
	Year string `json:"year"`
}