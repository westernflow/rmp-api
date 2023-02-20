package uwomodel

type Professor struct {
	Name           string   `json:"name" bson:"name"`
	RMPName        string   `json:"rmpName" bson:"rmpName"`
	RMPId          string   `json:"rmpId" bson:"rmpId"`
	Rating         float64  `json:"rating" bson:"rating"`
	Difficulty     float64  `json:"difficulty" bson:"difficulty"`
	CurrentCourses []string `json:"currentCourses" bson:"currentCourses"`
	Reviews        []Review `json:"reviews" bson:"reviews"`
	Departments    []string `json:"departments" bson:"departments"`
}

type Review struct {
	ProfessorID string  `json:"professorId" bson:"professorId"`
	Professor   string  `json:"professor" bson:"professor"`
	Quality     float64 `json:"quality" bson:"quality"`
	Difficulty  float64 `json:"difficulty" bson:"difficulty"`
	Date        string  `json:"date" bson:"date"`
	ReviewText  string  `json:"reviewText" bson:"reviewText"`
	Helpful     float64 `json:"helpful" bson:"helpful"` // quality = helpful+clarity/2
	Clarity     float64 `json:"clarity" bson:"clarity"`
}
