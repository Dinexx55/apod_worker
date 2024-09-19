package domain

type ApodImageMetaData struct {
	Id                    int    `json:"id" db:"id"`
	Title                 string `json:"title" db:"title"`
	Explanation           string `json:"explanation" db:"explanation"`
	Date                  string `json:"date" db:"date"`
	Copyright             string `json:"copyright" db:"copyright"`
	LocalStorageImagePath string `json:"localImagePath" db:"local_storage_path"`
}
