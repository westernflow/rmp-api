package model

type Review struct {
	Professor  string `json:"professor"`
	Quality 	int    `json:"quality"`
	Difficulty int    `json:"difficulty"`
	Date       string `json:"date"`
	ReviewText string `json:"reviewText"`
	Course     Course `json:"course"`
}

type Course struct {
	Department string `json:"department"`
	Number     int    `json:"number"`
	Reviews []Review `json:"reviews"`
}

// create a model for a professor with the fields: name, rating, numRatings, department, level of difficulty, and reviews
type Professor struct {
	Name        string  `json:"name"`
	Rating      float64 `json:"rating"`
	Department string `json:"department"`
	Difficulty  float64 `json:"difficulty"`
	Reviews     []Review `json:"reviews"`
}