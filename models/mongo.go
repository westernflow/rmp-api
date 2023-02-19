package model

type MongoProfessor struct {
	Name        string             `json:"name" bson:"name"`
	RMPId       string             `json:"rmpId" bson:"rmpId"`
	Rating      float64            `json:"rating" bson:"rating"`
	Difficulty  float64            `json:"difficulty" bson:"difficulty"`
	Reviews     []MongoReview      `json:"reviews" bson:"reviews"`
	Departments []MongoDepartment  `json:"departments" bson:"departments"`
}

type MongoReview struct {
	ProfessorID string             `json:"professorId" bson:"professorId"`
	Professor   string             `json:"professor" bson:"professor"`
	Quality     float64            `json:"quality" bson:"quality"`
	Difficulty  float64            `json:"difficulty" bson:"difficulty"`
	Date        string             `json:"date" bson:"date"`
	ReviewText  string             `json:"reviewText" bson:"reviewText"`
	Helpful     float64            `json:"helpful" bson:"helpful"` // quality = helpful+clarity/2
	Clarity     float64            `json:"clarity" bson:"clarity"`
}

type MongoDepartment struct {
	Name                 string             `json:"name" bson:"name"`
	DepartmentBase64Code string             `json:"departmentBase64Code" bson:"departmentBase64Code"`
	DepartmentNumber	 string             `json:"departNumber" bson:"departNumber"`
}
