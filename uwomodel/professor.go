package uwomodel

// Struct for a Western professor
type Professor struct {
	Name           string   `bson:"name" json:"name" examples:"J. Smith"`
	RMPName        string   `bson:"rmpName" json:"rmpName" examples:"John Smith"`
	RMPId          string   `bson:"rmpId" json:"rmpId" examples:"VGVhY2hlci0yMTM5MjE0"`
	Rating         float64  `bson:"rating" json:"rating" examples:"4.5"`
	Difficulty     float64  `bson:"difficulty" json:"difficulty" examples:"3.5"`
	CurrentCourses []string `bson:"currentCourses" json:"currentCourses" examples:"[\"COSC 101\", \"MATH 102\"]"`
	Reviews        []Review `bson:"reviews" json:"reviews"`
	Departments    []string `bson:"departments" json:"departments" examples:"[\"COSC\", \"MATH\"]"`
}

type Review struct {
	ProfessorID string `bson:"professorId" json:"professorId" examples:"VGVhY2hlci0yMTM5MjE0"`
	Quality     int    `bson:"quality" json:"quality" examples:"4"`
	Difficulty  int    `bson:"difficulty" json:"difficulty" examples:"5"`
	Date        string `bson:"date" json:"date" examples:"2018-02-07 10:28:42 +0000 UTC"`
	ReviewText  string `bson:"reviewText" json:"reviewText" examples:"This course was very interesting and the professor was very helpful."`
	Helpful     int    `bson:"helpful" json:"helpful" examples:"3"`
	Clarity     int    `bson:"clarity" json:"clarity" examples:"5"`
}
