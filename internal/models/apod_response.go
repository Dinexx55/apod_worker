package models

type APODResponse struct {
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	Date        string `json:"date"`
	Copyright   string `json:"copyright"`
	URL         string `json:"url"`
}
