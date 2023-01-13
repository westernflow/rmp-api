package model

type Review struct {
	Professor  string `json:"professor"`
	Quality 	int    `json:"quality"`
	Difficulty int    `json:"difficulty"`
	Date       string `json:"date"`
	ReviewText string `json:"reviewText"`
	Course     Course `json:"course"`
	Helpful 	 float64 `json:"helpful"` // quality = helpful+clarity/2	
	Clarity 	 float64 `json:"clarity"`
}

type Course struct {
	Department string `json:"department"`
	Number     int    `json:"number"`
}

// create a model for a professor with the fields: name, rating, numRatings, department, level of difficulty, and reviews
type Professor struct {
	Name        string  `json:"name"`
	RMPId 		 int     `json:"rmpId"`
	Rating      float64 `json:"rating"`
	Department string `json:"department"`
	Difficulty  float64 `json:"difficulty"`
	Reviews     []Review `json:"reviews"`
	Courses []Course `json:"courses"`
}

type Request struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}
